# simple-apns-server

    git clone https://github.com/manuks/simple-apns-server.git

    go build simple-apns-server

    chmod u+x simple-apns-server

    Usage of ./simple-apns-server:
      -cert="": The certificate file name
      -ip="127.0.0.1": The ip address that it lisents to
      -key="": The key file name
      -log=false: Show more log messages
      -port="8080": The port number
      -sandbox=false: Sandbox mode or not

  
    ./simple-apns-server -cert=pub.pem -key=key.pem -log -sandbox

    curl --data "message=message&badge=1&device=e9096b455bba8908f60e73c2dc57ede4f58002780174e35f84ccea0e48cb2674" http://localhost:8080/apn
