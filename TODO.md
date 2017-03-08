- [x] Implement `reader.NewDownstream()`
- [ ] Add support for collecting errors via `context.Context`
  - [ ] Errors can be collected anywhere in the `Context Tree`, and are propagated
  upwards until it reaches the root context.
  - [ ] Errors are signaled by `ctx.Close(err)` where `err != nil`.
  - [ ] Errors are stored in the root context. Need a way to allow any node in the
  tree to signal an error and propagate it to the root context.
  - [ ] Any node in the tree should be able to check for errors via `ctx.Err()`.
  - [ ] `rootCtx.Err()` returns any existing error, falling back to `goContext.Err()`
  in case the former is `nil`.
- [ ] Implement Producer abstraction
- [ ] Implement Transformer abstraction
- [ ] Implement Consumer abstraction
