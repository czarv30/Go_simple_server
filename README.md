# Go_simple_server
Simple training exercise.

Notes from phone conversation
- Write a server (in go) that connects to postgres database.
Post json blob, say student record. Put that into a student table. Maybe another classes table. 

Later notes (WhatsApp)
So to formalize the homework a bit more:
- install postgres
- create a table to house fields for student data
- in golang create a web server where I could execute:
-  a POST web requests against which will result in a student row being populated
- a GET web request which will result in all rows from that table being returned
- a different GET request where I could specify a couple of different parameters that will return to me just those rows that match
I should be able to recreate the same table on my end because you will have a .sql file that can be used to recreate it.

Other reqs
- Don't throw away errors.
- Define types to handle db data, and leverage the std library json parser for input/output. 
- Think about separating the SQL. 
- Submit updates via PRs. 
- Credentials CLI

Testing: 
- The Get request should be testable by pointing the browser to http://localhost:8080/GetStudents" 
- For the Post request I should be able to do:
curl -X POST -H "Content-Type: application/json" -d '{"first_name":"Dave","last_name":"Grohl","date_of_birth":"1971-04-25"}' http://localhost:8080/AddStudent