#!/bin/bash

set -e

if [ "$(pwd)" != "$GOPATH" ]; then
  echo "Setting GOPATH to the current directory"
  echo ""
  export GOPATH=$(pwd)
fi

export PATH="$GOPATH/bin:$PATH"

echo "Cleaning up build artifacts"
find . -type f -name '*.json' | xargs rm -f
rm -f datagen dsl/dsl.go

mkdir -p src/github.com/ThoughtWorksStudios

if [ ! -e src/github.com/ThoughtWorksStudios/datagen ]; then
  ln -s "$(pwd)" src/github.com/ThoughtWorksStudios/
fi

for dep in github.com/Pallinder/go-randomdata \
           github.com/mna/pigeon \
           ; do
  if [ ! -d "src/$dep" ]; then
    echo "Installing $dep..."
    go get -u "$dep"
  fi
done

echo ""
echo "Generating parser..."
pigeon -o dsl/dsl.go dsl/dsl.peg

echo ""
echo "Compiling..."
go build

echo ""
echo "processing example.lang"
./datagen example.lang

echo ""
echo "Running tests..."
go test github.com/ThoughtWorksStudios/datagen{,/dsl,/interpreter,/generator}

echo ""
echo "Done"
