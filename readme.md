### a simple server for slave browsers, like unattended wall-mounted displays. 

executing:

    ./run.sh

will "go get" some code from code.google.com, & start up server.

point browser(s) to http://localhost:8080 (default hard-coded url),
then access url's like this to set the image shown:

    curl -X PUT http://localhost:8080/post?url=http://www.clipartsfree.net/svg/18177-tv-test-screen-download.svg

or 

    curl -X PUT http://localhost:8080/post?url=http://www.clipartsfree.net/svg/tv-testscreen_Clipart_svg_File.svg

or

    curl -X PUT http://localhost:8080/post?url=https://www.google.com/images/srpr/logo11w.png
