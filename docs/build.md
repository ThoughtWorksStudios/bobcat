## Building from Source

The included Makefile has targets to get you started:

```bash
$ make list

Make targets:
   build clean compile depend local performance prepare release run smoke test wercker
```

Set up your dev workspace (this is meant to be ran the first time you set up your workspace). This will install golang from homebrew, configure the current directory for development, install dependencies, then finally build and run tests:

```bash
make local
```

Build and run tests:
```bash
make
```

Just build the binary:
```bash
make build
```

Just run tests:
```bash
make test
```

Install dependencies:
```bash
make depend
```
