Let's start with **Part 0.1**. Since you already have a strong programming background, I'll focus on the Python-specific mental model rather than explaining what a programming language is.

# 0.1 Python Development Environment

## 1. What Python actually is

Python is a **programming language and an interpreter-based runtime environment**.

When you write:

```python
name = "Harshad"
print(name)
```

you are writing Python source code. Your Mac cannot directly execute that source code as machine instructions.

Something needs to **read and execute the Python code**.

That's the job of the **Python interpreter**.

A simplified picture is:

```text
Your Python code
       ↓
Python interpreter
       ↓
Your program runs
```

This is different from the way you might commonly think about C#:

```text
C# source code
       ↓
C# compiler
       ↓
IL / assembly
       ↓
.NET runtime
       ↓
Program
```

Python doesn't normally have that same compile-to-IL workflow that you're familiar with from .NET.

Python does, however, involve an intermediate form called **bytecode**. We'll come back to that later; you don't need to understand it yet.

For now, the important idea is:

> **The Python interpreter is the program that knows how to execute Python code.**

---

# 2. Installing Python on macOS

Since I'm on macOS and prefer Homebrew, we'll use Homebrew.

First, check whether Homebrew is installed:

```bash
brew --version
```

You should get something similar to:

```text
Homebrew 4.x.x
```

If you don't have Homebrew installed, install it from the official Homebrew website.

Once Homebrew is available:

```bash
brew install python
```

Homebrew will install Python and the associated tools.

Now check it:

```bash
python3 --version
```

You should see something like:

```text
Python 3.x.x
```

The exact version will depend on the current Homebrew version when you install it.

---

# 3. Why `python3` instead of `python`?

This is one of those things that can initially be confusing.

You may see people type:

```bash
python
```

while others type:

```bash
python3
```

They are commands that can point to Python installations, but on macOS and other Unix-like systems, `python3` is commonly used explicitly for Python 3.

For this curriculum, we'll use:

```bash
python3
```

That makes it immediately obvious that we're running Python 3.

You can check which executable you're actually using:

```bash
which python3
```

You might get something like:

```text
/opt/homebrew/bin/python3
```

or another Homebrew-related path depending on your Mac.

This is useful because it answers an important question:

> **Which Python installation am I actually running?**

---

# 4. What is `pip3`?

Python itself isn't the only thing you'll need.

As you build applications, you'll install third-party packages.

For example, later we might install FastAPI:

```bash
pip3 install fastapi
```

`pip` is Python's package installer.

So conceptually:

```text
python3
   ↓
Runs Python programs

pip3
   ↓
Installs Python packages
```

You can check it with:

```bash
pip3 --version
```

You'll get something similar to:

```text
pip 25.x.x from ...
```

Notice that `pip3` is associated with a particular Python installation.

That's important later when we introduce **virtual environments**.

---

# 5. Python versions

You'll encounter versions such as:

```text
Python 3.10
Python 3.11
Python 3.12
Python 3.13
Python 3.14
```

Python's major version is currently **3**.

The other number represents the minor version.

For example:

```text
3.13.7
│ │  │
│ │  └── patch
│ └───── minor
└─────── major
```

You don't need to memorize Python's versioning rules yet.

What's important is that **Python versions can behave differently**, particularly when libraries depend on specific versions.

That's why professional Python projects explicitly define which Python versions they support.

---

# 6. Which Python version will we use?

For this curriculum, use a **current stable Python 3 release supported by the libraries we need**.

Don't blindly upgrade Python every time a new release appears.

For learning, consistency is more important than always having the newest version.

Once we start building AI applications, we'll also pay attention to what versions the libraries and model SDKs actually support.

After installing Python, run:

```bash
python3 --version
```

and keep that version in mind.

---

# 7. The Python interpreter

Here's something I want you to actually try.

Run:

```bash
python3
```

You should see something resembling:

```text
Python 3.x.x ...
>>>
```

That `>>>` is the **Python interactive prompt**.

You can now type Python directly:

```python
>>> 2 + 3
5
```

Try:

```python
>>> print("Hello, Python")
```

You should get:

```text
Hello, Python
```

This is called the **REPL**:

> **Read → Evaluate → Print → Loop**

Python:

1. Reads what you type.
2. Evaluates it.
3. Prints the result.
4. Waits for more input.

You can leave it with:

```python
>>> exit()
```

or press:

```text
Ctrl + D
```

on macOS.

---

# 8. Why is the REPL useful?

You won't normally build an application by typing everything into the REPL.

But it's incredibly useful for **quick experimentation**.

For example, later you might wonder:

> "What does this Python method actually return?"

Instead of creating an entire project, you can simply run:

```bash
python3
```

and experiment.

It's roughly analogous to having a quick interactive environment where you can test C# expressions, although Python's REPL is a much more central part of the everyday Python workflow.

---

# 9. Your first Python file

Now let's move from the interactive interpreter to an actual program.

Create a directory somewhere convenient:

```bash
mkdir python-playground
cd python-playground
```

Create a file:

```bash
touch hello.py
```

Open it in VS Code:

```bash
code .
```

If `code` isn't available in your terminal, you can open the folder through VS Code normally.

Put this in `hello.py`:

```python
print("Hello from Python!")
```

Now run:

```bash
python3 hello.py
```

You should see:

```text
Hello from Python!
```

And that's your first Python program.

Notice what happened:

```text
hello.py
   ↓
python3
   ↓
Python interpreter reads the file
   ↓
Program executes
```

The important distinction is:

```text
python3
```

is the **interpreter**.

```text
hello.py
```

is your **Python source code**.

---

# 10. One important distinction

You now have two ways to give Python code to the interpreter.

### Interactive

```bash
python3
```

Then:

```python
>>> print("Hello")
```

### Script

```bash
python3 hello.py
```

The first is useful for experimentation.

The second is how we'll normally build applications.

---

## Checkpoint

Before moving to **0.2**, you should be able to answer these in your own words:

1. What is the Python interpreter?
2. What is the difference between `python3` and `pip3`?
3. Why are we using `python3` rather than simply `python`?
4. What happens when you run `python3 hello.py`?
5. What is the `>>>` prompt?
6. What is the difference between using the REPL and running a `.py` file?

## Answers

1. The Python interpreter reads your Python source code and executes it.
2. `python3` runs Python programs; `pip3` installs Python packages.
3. `python3` means “use Python 3,” rather than Python 2. The important point is that `python3` explicitly selects the Python 3 interpreter.
4. `python3 hello.py` tells the Python 3 interpreter to read and execute `hello.py`.
5. `>>>` is the prompt shown by the Python REPL.
6. The REPL is useful for quickly experimenting with Python code without creating a `.py` file.

One useful distinction to keep in your head:

```text
python3 hello.py
       │
       └── "Python, execute this source file"
```

versus:

```text
python3
   ↓
>>>
   ↓
"Python, I'm going to give you code interactively."
```

### Practical task

On your Mac, do the following yourself:

```bash
python3 --version
pip3 --version
which python3
python3
```

Then run something in the REPL and create/run your own `hello.py`.

**Don't just copy the commands blindly.** The important part of this lesson is that you understand what each command is doing.

### Next Lesson → 0.2 Your First Python Program

We'll go one step further and look at what actually happens when a `.py` file is executed, including the relationship between your terminal, the interpreter, and your Python code.