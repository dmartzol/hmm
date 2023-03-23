# hurl --variable first_name="john" \
#      --variable last_name="doe" \
#      --variable email="john@example.com" \
#      --variable password="password12345" \
#      --variable dob="1990-01-01" \
#      account_create.hurl

POST http://localhost:1100/accounts
Content-Type: application/json
{

	"FirstName": "{{first_name}}",
	"LastName": "{{last_name}}",
	"Email": "{{email}}",
	"Password": "{{password}}",
	"DOB": "{{dob}}"
}

