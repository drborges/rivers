- [ ] Implement cancellation checks for `stream.Writer#Write`

- [ ] Implement `stream.NewDownstream(stream.Reader)`

  ```go
  func TestDownstreamCreation(t *testing.T) {
  	expect := expectations.New()

  	reader, writer := stream.New()

  	childReader, childWriter := stream.NewDownstream(reader)
  }
  ```

- [ ] Implement Producer abstraction

  ```go
  type Producer func(context.Context) stream.Reader
  ```

- [ ] Implement Transformer abstraction

  ```go
  type Transformer func(stream.Reader) stream.Reader
  ```

- [ ] Implement Consumer abstraction

  ```go
  type Consumer func(stream.Reader)
  ```
