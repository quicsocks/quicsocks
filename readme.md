1. go build
2. sudo mkdir /opt/quicsocks
3. sudo mv quicsocks /opt/quicsocks
4. sudo cp systemd/quicsocks.service /usr/lib/systemd/system/
5. sudo systemctl enable quicsocks
6. sudo systemctl start quicsocks
