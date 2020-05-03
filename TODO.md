# TODO


- [ ] Add Row fields in query db table
- [ ] validate emails with https://github.com/badoux/checkmail/blob/master/checkmail.go
- [ ] use PasswordStrength() to check password strength: https://github.com/nbutton23/zxcvbn-go
https://stackoverflow.com/questions/48345922/reference-password-validation
- [ ] https://ux.stackexchange.com/questions/110321/what-would-be-best-layout-for-registration-form-containing-14-input-fields
- [ ] Forgot password message: "If the account exists, an email will be sent with recovery details."
- [X] sessions: https://stackoverflow.com/questions/21680359/postgresql-create-access-token-on-insert/21684011#21684011
- [ ] job than runs CleanSessionsOlderThan() every week?
- [ ] return APIStructures(with no sensitive info) instead of original structures
- [ ] Parameters: https://stackoverflow.com/questions/4024271/rest-api-best-practices-where-to-put-parameters
- [ ] excluding from middleware: https://stackoverflow.com/questions/47957988/middleware-on-a-specific-route
- [ ] Reset password: https://stackoverflow.com/questions/3077229/restful-password-reset
- [ ] Testing: https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql
- [ ] w.WriteHeader(http.StatusAccepted) for all
- [ ] Better name for external ID's like: external_payment_customer_id