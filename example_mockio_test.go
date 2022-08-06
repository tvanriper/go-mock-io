package mockio_test

import (
	"fmt"
	"io"
	"strings"
	"time"

	mockio "github.com/tvanriper/go-mock-io"
)

// HandleUgliness could be a function you create to respond to serial input.
func HandleUgliness(b []byte) (count int, ok bool) {
	s := string(b)
	if strings.Contains(s, "ugly") {
		return len(b), true
	}
	return 0, false
}

// SendMessages could be your own function that you want to test.
func SendMessages(msg1 string, msg2 string, serial io.Writer) {
	serial.Write([]byte(msg1))
	time.Sleep(100 * time.Millisecond)
	serial.Write([]byte(msg2))
}

func ExampleMockIO() {
	// Setting up the mock serial port:
	serial := mockio.NewMockIO()
	msg1 := "Hullo, beautiful!"
	resp1 := "Hi, hansome!"
	serial.Expect(mockio.NewExpectBytes([]byte(msg1), []byte(resp1), 53*time.Millisecond))
	msg2 := "Woah, ugly!"
	resp2 := "Damn, nasty!"
	serial.Expect(mockio.NewExpectFunc(HandleUgliness, []byte(resp2), 43*time.Millisecond))

	// Setting up a reader thread.
	wait := make(chan struct{}, 1)
	go func() {
		for i := 0; i < 3; i++ {
			buffer := make([]byte, 100)
			n, err := serial.Read(buffer)
			if err != nil {
				fmt.Printf("error: %s\n", err)
			}
			fmt.Printf("received: \"%s\"\n", buffer[:n])
		}
		close(wait)
	}()

	// Use the serial port in your functions.
	SendMessages(msg1, msg2, serial)

	// Sending something unannounced to the serial port.
	time.Sleep(100 * time.Millisecond)
	serial.Send([]byte("Wait, don't go!"), 10*time.Millisecond)

	<-wait

	// Output:
	// received: "Hi, hansome!"
	// received: "Damn, nasty!"
	// received: "Wait, don't go!"
}
