# ctxlog

Ctxlog package provides a simple abstraction that allows to implicitly pass 
diagnostic information between layers of abstraction by adopting golang's context pattern.

[![GoDoc](https://godoc.org/github.com/andreiko/ctxlog?status.svg)](http://godoc.org/github.com/andreiko/ctxlog)

Example:
```
func (h *RequestHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
  ctx, logger := ctxlog.GetUpdatedLoggingContext(context.Background(), h.Logger, log.Fields{
    "request_id": new_uuid_v4(),
    "request_addr": request.RemoteAddr,
  })

  data, err := h.Controller.GetResponseData(ctx, request.URL.Path)
  if err != nil {
    logger.WithError(err).Error("error getting response data from controller")
    h.Error500(response, err)
    return
  }
  response.Write(data)
}

// ...

func (d *Controller) GetResponseData(ctx context.Context, path string) ([]byte, error) {
  return d.Storage.GetFile(ctx, path)
}

// ...

func (s *Storage) GetFile(ctx context.Context, path string) ([]byte, error) {
  logger := ctxlog.GetContextualLogger(ctx, s.Logger)
  data, err := s.FetchFile(ctx, path)
  if err != nil {
    // Message will contain request_id and request_addr fields from RequestHandler
    logger.WithError(err).Error("could not fetch file")
    return nil, err
  }

  return data, nil
}

```
