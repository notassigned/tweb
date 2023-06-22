cd ~/sandbox/tweb
if [ ! -f ./id ]; then
  go run . genkey > id2
fi
go run . ui.xml
cd -
