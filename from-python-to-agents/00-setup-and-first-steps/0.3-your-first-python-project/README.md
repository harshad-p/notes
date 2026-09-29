# 0.3 Your First Python Project

So far, we've been running Python directly with:

```bash
python3 hello.py
```

That works perfectly for a tiny experiment. But as soon as you have a real Python project, you'll need a way to keep **that project's Python environment and dependencies separate** from everything else on your machine.

That's what a **virtual environment** is for.

---

## 1. What problem are we solving?

Imagine you have two Python projects:

```text
Project A
  needs requests 2.x

Project B
  needs requests 3.x
```

If both projects use the same Python environment, their dependencies can interfere with each other.

You don't want:

> "I upgraded a package for Project B and suddenly Project A stopped working."

Python solves this with **virtual environments**.

Each project can have its own isolated environment containing its own installed packages.

---

# 2. Create a proper project directory

Let's create a new project rather than continuing to use our playground.

From somewhere appropriate, run:

```bash
mkdir my-first-python-project
cd my-first-python-project
```

We now have:

```text
my-first-python-project/
```

This is our project directory.

Eventually, a real project might look something like:

```text
my-first-python-project/
├── .venv/
├── src/
├── tests/
├── README.md
└── ...
```

We aren't going to create all of that yet.

For now, we'll focus on `.venv`.

---

# 3. Creating a virtual environment

Inside the project directory, run:

```bash
python3 -m venv .venv
```

There are two things worth understanding here.

### `python3`

We're using the Python interpreter we already installed.

### `-m venv`

`-m` tells Python:

> "Run a Python module as a program."

`venv` is Python's built-in module for creating virtual environments.

### `.venv`

This is the name we're giving our virtual environment.

So:

```bash
python3 -m venv .venv
```

essentially means:

> "Use this Python installation to create a virtual environment called `.venv`."

After running it, you'll have something like:

```text
my-first-python-project/
└── .venv/
```

The `.venv` directory contains the environment's Python executable and supporting files.

---

# 4. Why `.venv`?

You might wonder why we called it `.venv` instead of something like `environment`.

The dot makes it a hidden directory on macOS/Linux, and `.venv` has become a very common convention for Python projects.

You'll see this frequently in real Python repositories:

```text
project/
├── .venv/
├── src/
└── ...
```

We will also add `.venv/` to `.gitignore` later.

You **do not commit your virtual environment to Git**.

---

# 5. Activating the virtual environment

Creating the environment doesn't automatically make your terminal use it.

On macOS/Linux, run:

```bash
source .venv/bin/activate
```

Your terminal prompt should change slightly. You may see something like:

```text
(.venv) your-machine:my-first-python-project $
```

The important part is:

```text
(.venv)
```

That tells you the virtual environment is currently active.

---

# 6. What does "activate" actually mean?

This is an important concept.

When you activate `.venv`, your shell is configured so that commands such as:

```bash
python
pip
```

refer to the Python environment inside `.venv`.

Without activation, you might be using:

```text
/usr/local/.../python3
```

or another system/Homebrew Python installation.

With the virtual environment activated, you're effectively using:

```text
my-first-python-project/.venv/.../python
```

So the command:

```bash
python
```

now refers to the project's Python environment.

You can verify this:

```bash
which python
```

You should see a path containing:

```text
my-first-python-project/.venv/bin/python
```

This is one of the most useful commands when debugging Python environments.

---

# 7. Why not just use `python3` everywhere?

You can.

But once you're working inside a virtual environment, the usual workflow is:

```bash
python
pip
```

rather than:

```bash
python3
pip3
```

Why?

Because activation makes `python` and `pip` point to the project's virtual environment.

So:

```bash
python
```

means:

> "Use the Python interpreter belonging to this project."

And:

```bash
pip
```

means:

> "Install a package into this project's environment."

This becomes particularly important once different projects have different dependencies.

---

# 8. Installing a package

Let's install a small package.

We'll use `requests`, which is a popular Python library for making HTTP requests.

With `.venv` activated:

```bash
pip install requests
```

`pip` downloads the package and installs it into the current virtual environment.

You can verify it:

```bash
pip show requests
```

You can also see installed packages with:

```bash
pip list
```

You should see `requests` and its dependencies.

---

# 9. Why does this matter?

Suppose you have:

```text
Project A
└── .venv/
    └── requests

Project B
└── .venv/
    └── different packages
```

The `requests` installation in Project A belongs to **Project A's environment**.

Project B doesn't automatically get it.

That's the isolation we wanted.

This is very similar in spirit to how you might think about NuGet dependencies in .NET, although Python's environment/package model works differently.

---

# 10. Running code inside the virtual environment

Let's create:

```bash
touch hello.py
```

Put this inside:

```python
import requests

print("My Python project is running!")
print(requests.__version__)
```

Now run:

```bash
python hello.py
```

Because `.venv` is activated, Python uses the project's interpreter, which can see the `requests` package you installed.

You should get something similar to:

```text
My Python project is running!
2.x.x
```

The exact version isn't important here.

The important part is that Python found `requests` because it was installed into the active project's environment.

---

# 11. System Python vs project Python

This distinction is worth getting comfortable with.

### System/Homebrew Python

This is the Python installation available on your machine.

For example:

```bash
python3 --version
```

might show:

```text
Python 3.14.x
```

That's your general Python installation.

### Project Python

When `.venv` is activated:

```bash
python --version
```

you're using the interpreter belonging to the project.

It was created from your Python installation when you ran:

```bash
python3 -m venv .venv
```

So conceptually:

```text
Your machine
│
├── Python installation
│
├── Project A
│   └── .venv
│       └── project Python + packages
│
└── Project B
    └── .venv
        └── project Python + packages
```

The virtual environments give each project its own space for Python packages.

---

# 12. Deactivating the environment

When you're finished working on the project, run:

```bash
deactivate
```

The `(.venv)` should disappear from your terminal prompt.

You can check:

```bash
which python3
```

and you'll be back to your normal Python installation.

To work on the project again later:

```bash
cd my-first-python-project
source .venv/bin/activate
```

---

# 13. One thing to notice

You might be wondering:

> "If `.venv` contains Python, did we install Python again?"

Not quite in the way you might initially think.

We used your existing Python installation to **create an isolated environment** with its own interpreter/environment setup.

The key thing we're isolating here is the **project's Python environment and installed packages**.

You don't need to think of `.venv` as "another completely independent Python installation" for now. The important mental model is:

> **One machine-level Python installation can be used to create separate project environments.**

---

# The workflow to remember

For a new Python project, you'll commonly do:

```bash
mkdir my-project
cd my-project

python3 -m venv .venv
source .venv/bin/activate

pip install requests

python hello.py
```

And when you're done:

```bash
deactivate
```

That's the basic Python project workflow.

---

## Checkpoint

Do this yourself rather than just reading it:

### 1. Create the project

```bash
mkdir my-first-python-project
cd my-first-python-project
```

### 2. Create and activate the environment

```bash
python3 -m venv .venv
source .venv/bin/activate
```

### 3. Verify where Python comes from

```bash
which python
python --version
```

Look at the path returned by `which python`. It should point inside `.venv`.

### 4. Install `requests`

```bash
pip install requests
```

Then:

```bash
pip list
```

### 5. Create and run `hello.py`

```python
import requests

print("My project is running!")
print(requests.__version__)
```

Run:

```bash
python hello.py
```

### 6. Deactivate

```bash
deactivate
```

Then check:

```bash
which python3
```

---

### Make sure you can explain these four things

1. **Why do Python projects use virtual environments?**
2. **What does `python3 -m venv .venv` do?**
3. **What changes when you run `source .venv/bin/activate`?**
4. **Where does `pip install requests` install `requests` when `.venv` is active?**

Once those make sense, **0.3 is done**. Next we'll move to **0.4 Development Tools**, where we'll set up the practical workflow around VS Code, the Python extension, interpreter selection, Git, `.gitignore`, and a README.