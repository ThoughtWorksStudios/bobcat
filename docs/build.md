## Building from Source

The included Makefile has targets to get you started:

```
$ make list

Make targets:
   build clean compile depend local performance prepare release run smoke test wercker
```

Set up your dev workspace (this is meant to be ran the first time you set up your workspace). This will install golang from homebrew, configure the current directory for development, install dependencies, then finally build and run tests:

```
make local
```

Build and run tests:
```
make
```

Just build the binary:
```
make build
```

Just run tests:
```
make test
```

Install dependencies:
```
make depend
```
