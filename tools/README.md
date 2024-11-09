The order is important. model gen is a must; api is not in use at the moment however the web app does register api link thus it needs to be compilable as well.

```
go run tools/gen-model/main.go -type model
go run tools/gen-model/main.go -type api
```

Gen-model is independent so even all other go files is non-compilable this still runs. By fixing templates we can generate back the working code.

Form gen and the rest depend on the correctness of the above (model, not api). As of now I put it in the `tools/gen-form/main_test.go` for conviniency to click the test func to run if needed.
