I would **not** make this a “Python course followed by an AI course.” I'd make it a single progression:

**Python → Professional Python → AI foundations → LLM engineering → Agent engineering → Production AI → FDE/AI Engineer specialization**

And because I already know .NET, APIs, databases, Docker, Kubernetes, CI/CD, system design, etc., we can skip a lot of material that would otherwise be redundant.

# Python → AI Engineer / FDE Roadmap

## The end goal

By the end, I should be able to look at a problem and say:

> “I can build the Python service, integrate the LLM, give it tools, connect it to company data, make it reliable, observe it in production, and explain the architecture and trade-offs to a customer or engineering team.”

That is much closer to **AI Engineer / FDE** than simply knowing Python.

---

# Part 0 — Setup & First Steps

## 0.1 Python Development Environment

- What Python is and how it runs
- Installing Python on macOS using Homebrew
- Checking the Python installation
- `python3` and `pip3`
- Python versions
- Choosing the Python version for this curriculum
- Understanding the Python interpreter

## 0.2 Your First Python Program

- Creating a project folder
- Creating a `.py` file
- Writing the simplest Python program
- Running a Python file from the terminal
- Running Python interactively
- Understanding what happens when a `.py` file runs

## 0.3 Your First Python Project

- Creating a proper project directory
- Creating a virtual environment
- Activating and deactivating the environment
- Understanding why virtual environments exist
- Installing a package with `pip`
- Running code inside the virtual environment
- Understanding the difference between system Python and project Python

## 0.4 Development Tools

- VS Code setup for Python
- Python extension
- Selecting the correct interpreter
- Running and debugging Python
- Basic terminal workflow
- Git initialization
- `.gitignore`
- `README.md`

## 0.5 First Project Checkpoint

Build and run a tiny Python project from scratch.

You should be able to:

- Install Python
- Check the Python version
- Create a project directory
- Create and run a `.py` file
- Create and activate a virtual environment
- Install a package
- Run the project using the virtual environment
- Explain why the virtual environment exists
- Commit the project to Git

# Part 1 — Python Fundamentals

**Goal:** Become comfortable reading and writing Python without constantly translating everything mentally from C#.

### 1.1 Python mental model

- Python execution model
- Python interpreter
- Objects and references
- Dynamic typing
- Variables vs objects
- Mutability
- Identity vs equality
- `None`
- Truthiness
- Python's philosophy and conventions
- Python vs C# mental model

### 1.2 Basic syntax

- Numbers
- Strings
- Booleans
- Variables
- Operators
- Expressions
- Comments
- f-strings
- Basic input/output

### 1.3 Collections

- Lists
- Tuples
- Dictionaries
- Sets
- Indexing
- Slicing
- Membership
- Mutability
- Nested collections

### 1.4 Control flow

- `if / elif / else`
- `for`
- `while`
- `break`
- `continue`
- `range`
- Iteration
- Python's `for` model

### 1.5 Functions

- Defining functions
- Parameters
- Return values
- Positional arguments
- Keyword arguments
- Default arguments
- `*args`
- `**kwargs`
- Scope
- Closures

### 1.6 Pythonic data manipulation

- List comprehensions
- Dictionary comprehensions
- Set comprehensions
- Conditional expressions
- `enumerate`
- `zip`
- `any`
- `all`
- `map`
- `filter`

### 1.7 Modules and packages

- Imports
- Modules
- Packages
- `__name__`
- `__main__`
- Python package structure
- Virtual environments
- `pip`

### 1.8 Errors and exceptions

- Exceptions
- `try`
- `except`
- `else`
- `finally`
- Raising exceptions
- Custom exceptions
- Exception design

### 1.9 Files and resources

- Reading/writing files
- Paths
- `pathlib`
- Context managers
- `with`

### 1.10 Classes and objects

- Classes
- Constructors
- Instance attributes
- Methods
- Class attributes
- Properties
- Inheritance
- Composition
- Special methods / dunder methods

### 1.11 Dataclasses and data modeling

- `dataclass`
- Immutable/frozen dataclasses
- Enums
- Data validation
- When to use classes vs dictionaries

### 1.12 Iterators and generators

- Iterables
- Iterators
- `yield`
- Generators
- Lazy evaluation

This is particularly important later for **streaming and data processing**.

---

# Part 2 — Professional Python

Now we move from:

> “I can write Python”

to:

> “I can maintain a Python codebase.”

### 2.1 Type hints

- Type annotations
- `list[str]`
- `dict[str, int]`
- `Optional`
- Union types
- `Any`
- `Literal`
- Generics
- Type aliases
- Static type checking

### 2.2 Protocols and abstraction

- Duck typing
- Protocols
- Abstract base classes
- Structural typing
- Dependency inversion in Python

We'll explicitly compare these with **C# interfaces**.

### 2.3 Decorators

- What decorators actually do
- Function wrapping
- `@decorator`
- Parameterized decorators
- Practical uses
- Why frameworks use them

This becomes important later because you'll encounter decorators everywhere in Python AI frameworks.

### 2.4 Context managers

- `with`
- `__enter__`
- `__exit__`
- `contextlib`
- Resource management

### 2.5 Project structure

Build a real Python project:

```text
project/
├── src/
├── tests/
├── pyproject.toml
├── README.md
└── ...
```

### 2.6 Dependency management

- `pyproject.toml`
- Dependency specification
- Virtual environments
- Lock files
- Development dependencies
- Reproducible environments

### 2.7 Testing

- `pytest`
- Fixtures
- Parametrized tests
- Mocking
- Integration tests
- Test doubles
- Testing async code

### 2.8 Logging and configuration

- `logging`
- Structured logging
- Environment variables
- Configuration management
- Secrets

### 2.9 Code quality

- Formatting
- Linting
- Static analysis
- `ruff`
- Type checking
- Code organization

### 2.10 Debugging and profiling

- Debugging Python
- Stack traces
- IDE debugging
- Profiling
- Finding bottlenecks
- Memory considerations

---

# Part 3 — Python Backend Engineering

You already know backend engineering, so this section is deliberately shorter.

The objective is to learn **how backend engineering looks in Python**.

### 3.1 HTTP in Python

- HTTP clients
- Requests
- Async HTTP
- JSON
- Headers
- Authentication
- Timeouts
- Retries

### 3.2 FastAPI

- Routes
- Request models
- Response models
- Dependency injection
- Validation
- Error handling
- Middleware
- Async endpoints
- OpenAPI

### 3.3 Database access

- PostgreSQL from Python
- Drivers
- SQLAlchemy
- ORM vs SQL
- Transactions
- Connection pooling
- Async database access

### 3.4 Background processing

- Background jobs
- Queues
- Workers
- Retries
- Idempotency
- Dead-letter handling

### 3.5 Production API

Build a real service with:

```text
FastAPI
   ↓
Service layer
   ↓
Repository
   ↓
PostgreSQL
```

with:

- authentication
- logging
- tests
- configuration
- Docker
- CI

At this point, Python becomes a **professional backend language** for you.

---

# Part 4 — AI Foundations

Now we start changing direction.

The goal isn't to become a research ML scientist.

The goal is:

> **Understand enough AI to make good engineering decisions.**

### 4.1 Machine learning mental model

- Training
- Inference
- Features
- Labels
- Models
- Parameters
- Loss
- Optimization
- Generalization

### 4.2 Neural networks

Conceptually understand:

- Neurons
- Layers
- Weights
- Activations
- Forward pass
- Backpropagation
- Training

Not a giant math detour.

### 4.3 Transformers

This is essential.

Understand:

- Tokens
- Embeddings
- Attention
- Self-attention
- Transformer architecture
- Encoder vs decoder
- Why transformers changed NLP

### 4.4 LLMs

- Pretraining
- Instruction tuning
- Fine-tuning
- RLHF / preference optimization at a conceptual level
- Inference
- Context windows
- Parameters
- Quantization
- Model sizes

### 4.5 Embeddings

- What embeddings represent
- Similarity
- Vector space
- Semantic search
- Embedding models

### 4.6 AI limitations

- Hallucination
- Context limitations
- Non-determinism
- Bias
- Data leakage
- Model uncertainty

This section prevents you from treating LLMs like magical APIs.

---

# Part 5 — LLM Application Engineering

This is where things get much more practical.

### 5.1 Calling LLM APIs

- API authentication
- Requests
- Responses
- Messages
- Streaming
- Tokens
- Usage
- Costs

### 5.2 Prompt engineering

Not “50 magic prompts.”

Instead:

- System instructions
- User instructions
- Context
- Few-shot examples
- Constraints
- Output specifications
- Prompt composition
- Prompt versioning

### 5.3 Structured outputs

- JSON
- Schemas
- Typed responses
- Validation
- Handling malformed outputs

### 5.4 Tool/function calling

Understand exactly what happens:

```text
User
 ↓
LLM
 ↓
"I need tool X"
 ↓
Application
 ↓
Tool execution
 ↓
Result
 ↓
LLM
 ↓
Final response
```

This becomes the foundation for agents.

### 5.5 Streaming

- Token streaming
- Server-sent events
- Async processing
- Partial responses

### 5.6 Context engineering

- Context windows
- Context selection
- Context compression
- Conversation history
- Summarization
- Relevant context retrieval

---

# Part 6 — RAG

RAG is important enough to get its own section.

### 6.1 Why RAG exists

Understand the fundamental problem:

> The model doesn't automatically know your company's private/current information.

### 6.2 Document ingestion

```text
Documents
   ↓
Parsing
   ↓
Chunking
   ↓
Embeddings
   ↓
Vector database
```

### 6.3 Retrieval

- Vector search
- Similarity
- Metadata filters
- Keyword search
- Hybrid search
- Reranking

### 6.4 Generation

```text
Question
 ↓
Retrieve relevant information
 ↓
Construct context
 ↓
LLM
 ↓
Answer
```

### 6.5 Production RAG

- Chunking strategies
- Retrieval quality
- Evaluation
- Citation
- Permissions
- Document freshness
- Incremental ingestion
- Caching

### Project

Build a **production-style knowledge assistant** over a collection of documents.

---

# Part 7 — Agent Engineering

This becomes one of the core parts of your curriculum.

## 7.1 What is an agent?

First distinguish:

**LLM application**

vs

**workflow**

vs

**agent**

vs

**multi-agent system**

This distinction is extremely important.

---

## 7.2 Tool-using agents

Build an agent that can use:

- APIs
- Database queries
- Search
- Calculator
- File operations

Understand:

- Tool definitions
- Tool schemas
- Tool selection
- Tool execution
- Tool results

---

## 7.3 Agent loops

Understand the loop:

```text
Observe
   ↓
Reason
   ↓
Choose action
   ↓
Execute tool
   ↓
Observe result
   ↓
Repeat
```

Then understand why blindly allowing this loop can be dangerous.

---

## 7.4 Agent state

- Conversation state
- Task state
- Intermediate results
- Checkpoints
- State persistence

---

## 7.5 Memory

- Short-term memory
- Long-term memory
- Semantic memory
- Episodic memory
- Memory retrieval
- Memory corruption
- When memory is unnecessary

---

## 7.6 Agent orchestration

- Sequential workflows
- Parallel workflows
- Conditional workflows
- Human approval
- Retry loops
- State machines

---

## 7.7 MCP

Learn:

- What MCP is
- Why it exists
- MCP clients
- MCP servers
- Tools
- Resources
- Prompts
- Security considerations

---

## 7.8 Multi-agent systems

Only after mastering single-agent systems:

- Specialist agents
- Supervisor agents
- Delegation
- Agent communication
- Shared state
- Coordination problems

And importantly:

> When **not** to use multi-agent architecture.

---

# Part 8 — Agent Reliability

This is where you start becoming an **AI Engineer rather than an AI hobbyist**.

### 8.1 Guardrails

- Input validation
- Output validation
- Tool restrictions
- Permissions
- Sandboxing

### 8.2 Prompt injection

- Direct injection
- Indirect injection
- Tool manipulation
- Data exfiltration
- Defense strategies

### 8.3 Reliability

- Retries
- Timeouts
- Idempotency
- Fallbacks
- Model fallback
- Tool failure
- Partial failure

### 8.4 Evaluation

This is critical.

Learn:

- Unit testing agents
- Golden datasets
- LLM-as-judge
- Retrieval evaluation
- Tool-use evaluation
- Regression testing
- End-to-end evaluation

### 8.5 Observability

- Logs
- Metrics
- Traces
- Token usage
- Latency
- Tool calls
- Agent trajectories
- Cost

---

# Part 9 — Production AI Engineering

Now combine your existing backend knowledge with everything above.

### Architecture

Build systems involving:

```text
Client
   ↓
API
   ↓
Application
   ↓
Agent / Workflow
   ↓
LLM
   ↓
Tools
   ↓
External systems
```

with:

```text
        ┌── PostgreSQL
        ├── Redis
Agent ──┼── Vector DB
        ├── External APIs
        └── Queues
```

### Topics

- Authentication
- Authorization
- Secrets
- Rate limiting
- Caching
- Queues
- Workers
- Async processing
- Cost controls
- Model routing
- Observability
- Evaluation
- Security
- Data privacy
- PII
- GDPR considerations
- Docker
- Cloud deployment
- CI/CD
- Production incident handling

Your existing **Docker/Kubernetes/CI/CD/system-design knowledge** becomes very valuable here, so we won't reteach it from scratch.

---

# Part 10 — FDE Track

This is where the curriculum branches.

The technical foundation remains shared with AI Engineering.

But FDE adds another dimension:

> **Solve real customer problems with technology.**

### 10.1 Customer discovery

- Asking the right questions
- Finding the actual problem
- Separating symptoms from requirements
- Technical discovery

### 10.2 Rapid prototyping

Given:

> “Our support team spends 4 hours every day doing X.”

You should be able to quickly determine:

- Can AI help?
- What data is required?
- What integrations are required?
- What architecture makes sense?
- What should be automated?
- Where should a human remain involved?

### 10.3 Integration engineering

- APIs
- Webhooks
- OAuth
- Authentication
- Data mapping
- ETL
- Customer systems
- Legacy systems

### 10.4 Customer environments

- Deployment constraints
- Security requirements
- Networking
- Data residency
- Permissions
- Enterprise authentication

### 10.5 Communication

Learn to explain:

> “Here's what we built.”

to:

- Engineers
- Product managers
- Executives
- Customers
- Nontechnical stakeholders

---

# Part 11 — AI Engineer Track

For the AI Engineer specialization, go deeper into:

### AI architecture

- Model selection
- Model routing
- RAG architecture
- Agent architecture
- Evaluation architecture

### AI infrastructure

- Inference
- GPU basics
- Model serving
- Quantization
- Local models
- Ollama
- vLLM
- Model APIs

### Advanced LLM engineering

- Fine-tuning
- LoRA / PEFT
- Function calling
- Structured generation
- Advanced retrieval
- Agent evaluation
- Synthetic data

You don't necessarily need all of this immediately for an AI Engineer role, so these come **after** the core.

---

# Part 12 — Portfolio

I'd make the projects progressively resemble real engineering work.

### Project 1 — Python Developer Tool

A useful CLI application.

**Goal:** Python fundamentals.

---

### Project 2 — Production Python API

FastAPI + PostgreSQL + tests + Docker + CI.

**Goal:** Professional Python.

---

### Project 3 — External API Integration Platform

Multiple external APIs, authentication, retries, caching and background jobs.

**Goal:** Integration engineering.

This is particularly relevant to FDE.

---

### Project 4 — LLM Application

A genuinely useful LLM application with:

- structured output
- streaming
- persistence
- evaluation

---

### Project 5 — RAG System

A production-style knowledge system with:

- ingestion
- chunking
- embeddings
- vector search
- metadata filtering
- citations
- evaluation

---

### Project 6 — Tool-Using Agent

An agent capable of:

- searching
- calling APIs
- querying data
- executing tools
- maintaining state

---

### Project 7 — Production Agent Platform

This is the serious one.

Something like:

```text
                    ┌── Web Search
                    ├── PostgreSQL
User → API → Agent ─┼── Internal APIs
                    ├── Knowledge Base
                    └── External Services
```

with:

- authentication
- tool permissions
- memory
- RAG
- evaluation
- tracing
- retries
- observability
- cost tracking
- human approval

This becomes a strong portfolio piece.

---

### Project 8 — FDE Simulation

You receive a fictional customer problem.

For example:

> A large company wants an AI assistant that can answer questions about internal operational data and take actions through their existing systems.

You have to:

1. Discover requirements
2. Identify constraints
3. Design architecture
4. Build a prototype
5. Integrate APIs
6. Handle authentication
7. Demonstrate it
8. Explain trade-offs
9. Identify production risks

This is essentially **FDE practice**.

---

# Part 13 — Interview Preparation

We'll eventually create dedicated tracks for:

### Python

- Python coding
- Debugging
- Async
- Data structures
- Testing

### Backend

- API design
- Databases
- Distributed systems
- Queues
- Caching

### AI

- LLM architecture
- RAG
- Embeddings
- Evaluation
- Agents

### System design

Examples:

> Design a production RAG platform.

> Design an AI customer-support agent.

> Design an agent that can execute financial operations safely.

### FDE

Examples:

> A customer wants an AI agent but cannot give you direct database access. What do you do?

> A customer's API keeps timing out. How would you investigate?

> A customer wants a prototype in two days. How do you decide what to build?

---

# What we will deliberately NOT do

Initially, we will **not** prioritize:

- Competitive-programming-style Python
- Advanced Python internals
- Django
- GUI development
- Data science notebooks
- Pandas-heavy data analysis
- Traditional ML algorithms in depth
- Computer vision
- NLP before transformers
- Training neural networks from scratch
- Mathematics-heavy ML

Some of these can be added later if a target job requires them.

---

# The progression

The whole journey looks roughly like this:

```text
                    ┌────────────────────┐
                    │ Python Fundamentals│
                    └─────────┬──────────┘
                              ↓
                    ┌────────────────────┐
                    │ Professional Python│
                    └─────────┬──────────┘
                              ↓
                    ┌────────────────────┐
                    │ Python Backend     │
                    └─────────┬──────────┘
                              ↓
                    ┌────────────────────┐
                    │ AI Foundations     │
                    └─────────┬──────────┘
                              ↓
                    ┌────────────────────┐
                    │ LLM Engineering    │
                    └─────────┬──────────┘
                              ↓
                    ┌────────────────────┐
                    │ RAG                │
                    └─────────┬──────────┘
                              ↓
                    ┌────────────────────┐
                    │ Agent Engineering │
                    └─────────┬──────────┘
                              ↓
                    ┌────────────────────┐
                    │ Production AI     │
                    └─────────┬──────────┘
                              ↓
                 ┌────────────┴────────────┐
                 ↓                         ↓
          ┌──────────────┐          ┌──────────────┐
          │ AI Engineer  │          │     FDE      │
          │ specialization│          │ specialization│
          └──────────────┘          └──────────────┘
```

### The skill stack we're ultimately building

```text
                 AI / Agent Architecture
                         ▲
                         │
                 LLM Engineering
                         ▲
                         │
                 Python Engineering
                         ▲
                         │
              Backend / System Design
                         ▲
                         │
              Your existing experience
```
