package bleflows

import (
	"context"
	"errors"
	"testing"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/stretchr/testify/require"
)

// feed returns a channel preloaded with one packet per response; each packet is the
// response's index, which fakeDecode maps back.
func feed(responses []blecommands.Response) (<-chan []byte, func([]byte) (blecommands.Response, error)) {
	ch := make(chan []byte, len(responses))
	for i := range responses {
		ch <- []byte{byte(i)}
	}
	return ch, func(buf []byte) (blecommands.Response, error) {
		if r, ok := responses[buf[0]].(errResponse); ok {
			return nil, r.err
		}
		return responses[buf[0]], nil
	}
}

type errResponse struct {
	blecommands.Response
	err error
}

func complete() *blecommands.Status {
	return &blecommands.Status{Status: blecommands.StatusComplete}
}

func TestReceiveUntilCompleteCollectsTypedResponses(t *testing.T) {
	ch, decode := feed([]blecommands.Response{
		&blecommands.KeypadCodeCount{Count: 2},
		&blecommands.Status{Status: blecommands.StatusAccepted},
		&blecommands.KeypadCode{CodeID: 1, Name: "a"},
		&blecommands.KeypadCode{CodeID: 2, Name: "b"},
		complete(),
		&blecommands.KeypadCode{CodeID: 3, Name: "after complete"},
	})

	var codes []blecommands.KeypadCode
	err := receiveUntilComplete(context.Background(), ch, decode, collectInto(&codes))

	require.NoError(t, err)
	require.Len(t, codes, 2)
	require.Equal(t, uint16(1), codes[0].CodeID)
	require.Equal(t, uint16(2), codes[1].CodeID)
}

func TestReceiveUntilCompleteNilHandler(t *testing.T) {
	ch, decode := feed([]blecommands.Response{
		&blecommands.Status{Status: blecommands.StatusAccepted},
		complete(),
	})
	require.NoError(t, receiveUntilComplete(context.Background(), ch, decode, nil))
}

func TestReceiveUntilCompleteHandlerEndsEarly(t *testing.T) {
	// No StatusComplete follows the ID, as with AddKeypadCode.
	ch, decode := feed([]blecommands.Response{
		&blecommands.KeypadCodeID{CodeID: 7},
	})

	var id uint16
	err := receiveUntilComplete(context.Background(), ch, decode, func(res blecommands.Response) bool {
		r, ok := res.(*blecommands.KeypadCodeID)
		if ok {
			id = r.CodeID
		}
		return ok
	})
	require.NoError(t, err)
	require.Equal(t, uint16(7), id)
}

func TestReceiveUntilCompleteDecodeError(t *testing.T) {
	wantErr := errors.New("device error")
	ch, decode := feed([]blecommands.Response{
		&blecommands.KeypadCode{CodeID: 1},
		errResponse{err: wantErr},
		complete(),
	})

	var codes []blecommands.KeypadCode
	err := receiveUntilComplete(context.Background(), ch, decode, collectInto(&codes))
	require.ErrorIs(t, err, wantErr)
}

func TestReceiveUntilCompleteContextDone(t *testing.T) {
	ch, decode := feed([]blecommands.Response{
		&blecommands.KeypadCode{CodeID: 1},
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := receiveUntilComplete(ctx, ch, decode, nil)
	require.ErrorIs(t, err, context.Canceled)
}
