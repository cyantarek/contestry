Machine Setup: Production

01. Set Env values in the .env file
02. Install Mysql (linux: sudo apt-get install mysql-server, windows: install xampp)
02. Create Mysql Database "contest"
02. Install Docker / or sudo apt-get install lxc lxd -y
03. Docker pull tarek5/sandbox:v2
04. Run the single binary


TODO

01. Repeat database entry validation
02. Fake JWT with changed user_type validation on JWTValidator() to check if the username is really that type
03. Role check before go to meat of the controller