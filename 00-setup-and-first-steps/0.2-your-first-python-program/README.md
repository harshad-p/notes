# 0.2 Your First Python Program

You've already seen the basic idea of running Python. Now let's actually build a tiny Python program and, more importantly, understand **what is happening when you run it**.

We won't introduce variables, functions, or other Python concepts yet. The goal here is simply to understand the basic workflow.

---

## 1. Creating a project folder

Let's create a small folder for our Python experiments.

In your terminal:

```bash
mkdir python-playground
cd python-playground
```

You now have:

```text
python-playground/
```

Think of this as your workspace. Later, your projects will have more structure, but for now we deliberately keep it simple.

---

## 2. Creating a `.py` file

Create a file called:

```text
hello.py
```

You can do this from the terminal:

```bash
touch hello.py
```

Your folder now looks like:

```text
python-playground/
└── hello.py
```

The `.py` extension tells us that this is a **Python source file**.

Just like:

```text
Program.cs
```

is a C# source file, `hello.py` is a Python source file.

---

## 3. Your first Python program

Open `hello.py` and put this inside:

```python
print("Hello, Python!")
```

That's a complete Python program.

`print()` is a built-in Python function that displays something as output.

Don't worry about the details of functions yet. For now, just understand:

```python
print("Hello, Python!")
```

means:

> "Python, display this text."

---

## 4. Running the Python file

From inside `python-playground`, run:

```bash
python3 hello.py
```

You should see:

```text
Hello, Python!
```

Notice what you actually typed:

```text
python3 hello.py
```

There are two important pieces here:

```text
python3        hello.py
   │              │
   │              └── the source code you want to execute
   │
   └── the Python interpreter
```

You're essentially telling your computer:

> "Use the Python 3 interpreter to execute the code in `hello.py`."

This is the same basic idea as:

```bash
dotnet run
```

in a .NET project, although the underlying execution model is different.

---

## 5. What actually happens?

This is the important part of this lesson.

When you run:

```bash
python3 hello.py
```

the shell first finds the `python3` executable.

Then Python starts and reads the source code in `hello.py`.

It encounters:

```python
print("Hello, Python!")
```

Python executes that instruction, which produces:

```text
Hello, Python!
```

Then the program finishes.

Conceptually:

```text
Terminal
   │
   │ python3 hello.py
   ▼
Python interpreter
   │
   │ reads hello.py
   ▼
print("Hello, Python!")
   │
   ▼
Hello, Python!
```

That's the basic Python execution loop you'll be working with throughout this curriculum.

---

# 6. Running Python interactively

There is another way to use Python that you already saw in the previous lesson.

Just run:

```bash
python3
```

You should get something similar to:

```text
Python 3.x.x
>>>
```

Now type:

```python
print("Hello, Python!")
```

You'll immediately get:

```text
Hello, Python!
```

You didn't create a `.py` file this time.

You gave the code directly to the interpreter.

For example:

```text
>>> 2 + 3
5

>>> print("Hello")
Hello

>>> 10 * 5
50
```

This is the **REPL** you learned about in 0.1.

---

## 7. Script vs REPL

So there are two basic ways you'll interact with Python:

### Python script

You write code in a file:

```text
hello.py
```

and execute it:

```bash
python3 hello.py
```

This is what you'll primarily use when building actual applications.

### Python REPL

You start Python:

```bash
python3
```

and experiment interactively:

```text
>>> 2 + 3
5
```

The REPL is particularly useful for quickly testing an idea.

For example, if you're unsure what a particular Python operation does, you can often test it in a few seconds without creating a file.

---

## 8. One subtle but important thing

When you run:

```bash
python3 hello.py
```

the `.py` file itself isn't the thing "running."

Your **Python interpreter is running**, and it is reading and executing the instructions contained in `hello.py`.

That's a useful mental model:

> **The `.py` file contains the instructions. The Python interpreter executes them.**

We'll build on this idea later when we discuss modules, imports, packages, virtual environments, and eventually how Python applications are structured.

---

# Checkpoint

Before moving to **0.3**, make sure you can do this yourself.

### Terminal

```bash
mkdir python-playground
cd python-playground
touch hello.py
```

Put this in `hello.py`:

```python
print("Hello, Python!")
```

Then run:

```bash
python3 hello.py
```

You should see:

```text
Hello, Python!
```

Then try the REPL:

```bash
python3
```

and:

```python
>>> 2 + 3
5
>>> print("Testing Python")
Testing Python
```

Exit with:

```python
exit()
```

### And explain these in your own words:

1. What is the difference between `hello.py` and `python3`?
2. What happens when you run `python3 hello.py`?
3. What's the difference between running a `.py` file and using the REPL?
4. Why would you use the REPL if you can write `.py` files?

Once you've done that, we'll move to **0.3 Your First Python Project**, where we'll introduce the **virtual environment (`venv`)** and why Python projects normally shouldn't use the system Python directly.