if [ ! -x "test.sh" ]; then
  echo "test.sh is not executable; run: sudo chmod +x test.sh" >&2
  exit 1
fi
./test.sh

if [ $? -eq 0 ]; then
    echo "\nAll tests passed. Proceeding to build..."
else
    echo "\nTests failed. Aborting build..." >&2
    exit 1
fi

go build

if [ $? -eq 0 ]; then
    echo "Build successful."
else
    echo "Build failed. Please check the errors above." >&2
fi