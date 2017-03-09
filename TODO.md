- [x] Implement `reader.NewDownstream()`
- [ ] Add support for collecting errors via `context.Context`
  - [x] Closing the context with error causes the whole context tree to be closed
  - [x] Errors can be signaled anywhere in the `Context Tree`, and propagated
  upwards until it reaches the root context.
  - [x] Errors are signaled by `ctx.Close(err)` where `err != nil`.
  - [x] Errors are stored in the root context. Need a way to allow any node in the
  tree to signal an error and propagate it to the root context.
  - [ ] `rootCtx.Err()` returns any existing error, falling back to `goContext.Err()`
  in case the former is `nil`.
- [ ] Implement Producer abstraction
- [ ] Implement Transformer abstraction
- [ ] Implement Consumer abstraction
