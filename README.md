# go-mock-io
Test your socket i/o in pure golang

I needed a library to help me test working with a serial connection.  Unfortunately, I didn't see one readily available in Golang that served my needs.  So, I figured, I'll have to make one.

## Basic concept

Create a mock io.ReadWriteCloser, and give it certain expectations.  When someone writing to the mock stream meets any of those expectations, the stream responds with that expectation's response.

