if [ ! -f ./id ]; then
  echo "generating new id file"
  go run . genkey > id
else
  echo "id file already exists"
fi