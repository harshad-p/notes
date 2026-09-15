# 1.1 What Is a Distributed System?

A distributed system is a system whose work is performed by multiple independent computers, called **nodes**, that communicate over a network to provide one service to its users. A node can be a server, a virtual machine, a container, or a managed database instance. From a user's perspective, a request to `api.example.com` may look like one interaction; internally, several nodes may participate.

## Start with one machine

A single-machine system keeps the application code, its running process, and its data on one computer. For a small application, this model is appealing: requests enter one process, the process reads or writes local data, and there is no network communication between internal components. Debugging and reasoning are comparatively simple because the machine either works or it does not.

The limitation is that one machine has finite CPU, memory, disk capacity, network bandwidth, and a finite number of requests it can process in a given time. It is also a single point of failure: if that machine becomes unavailable, so does the service.

Consider an API that stores photos. Initially, one server accepts uploads, generates thumbnails, and stores metadata. As usage grows, thumbnail generation consumes CPU, uploads consume network capacity, and a hardware or deployment failure makes the entire application unavailable. Moving those responsibilities onto separate machines changes the system from a single-machine system into a distributed one.

```text
Client
  |
  v
API server  --->  database server
  |
  v
thumbnail worker
```

Each box is independently running software on a node. The API server now depends on the database server being reachable, and the worker receives work across the network rather than through an in-process function call. Those dependencies are what make the system distributed.

## Why distribute a system?

The most common reason is **scalability**: the ability of a system to handle more work as demand grows. If a service needs more request-processing capacity, adding another API server can increase the total capacity. If it needs more storage or database throughput, the data layer may eventually need to be expanded as well. Distribution lets us use resources from multiple machines rather than being limited to one.

Distribution is also used for **availability**, meaning the likelihood that the service is usable when a user needs it. Running two API servers means one can keep serving traffic if the other crashes. Availability is about continued service, not about preserving every piece of data. We will later distinguish it carefully from durability.

**Reliability** is the broader property of delivering correct service consistently over time. An available API that sometimes returns the wrong balance is not reliable. Redundant machines can help reliability, but they can also introduce disagreements between copies of data, so redundancy alone does not guarantee it.

**Performance** describes how efficiently the system responds, commonly in terms of latency and throughput. **Latency** is the time for one operation to complete. **Throughput** is the amount of work completed per unit of time, such as requests per second. Distributing a CPU-heavy thumbnail task to workers may improve API latency and total throughput, but adding network calls can also add latency. Distribution is therefore a trade-off, not a free performance upgrade.

**Fault tolerance** is the ability to continue providing an acceptable level of service when some components fail. For example, if one of three API servers crashes and the other two continue serving requests, the API tier tolerates that server failure. Fault tolerance is a design property that supports availability.

**Durability** means that acknowledged data survives failures. If the photo API tells a user that an upload succeeded, the photo and its metadata should still exist after a process restart, a machine failure, or an outage within the failure scope the system promised to handle. Availability asks, “Can I use the service now?” Durability asks, “Will accepted data still be there later?” A system can be temporarily unavailable yet durable, or available while risking loss of recently accepted data.

## What becomes difficult?

A local function call either returns a result or raises an error in the same process. A network request has more ambiguous outcomes. The remote node may be slow, unreachable, crashed after doing the work, or the response may have been lost on the network. When an API call times out, the caller often cannot tell whether the other side did nothing or completed the operation just before the reply was lost.

Machines also fail independently. One API server may be healthy while its database connection is broken; one copy of data may be updated while another is temporarily behind. Network communication is slower and less predictable than memory access, and adding machines introduces configuration, deployment, monitoring, and security concerns.

This is the central tension of distributed systems: we distribute work to gain capacity and reduce single-machine failures, but coordination across machines introduces new failure modes and uncertainty. Good system design makes those trade-offs explicit. It starts with the simplest design that satisfies the actual requirements, then distributes a component only when a concrete limitation or reliability need justifies it.

## Recap

- A distributed system is one service implemented by multiple networked nodes.
- We distribute systems for scalability, availability, fault tolerance, storage capacity, and sometimes performance.
- Scalability is handling more work; availability is being usable; reliability is consistently correct behavior; durability is survival of acknowledged data; fault tolerance is continued acceptable service after failures.
- Distribution removes some single-machine limits but adds network uncertainty, independently failing components, and coordination problems.

## Check your understanding

A single API server stores order data on its local disk. It is replaced by two API servers, each still storing orders only on its own local disk. A request may reach either server. What correctness problem can occur when a customer creates an order and then immediately fetches it, and why is this a distributed-systems problem rather than merely a scaling change?

---

Answer: If the fetch reaches the other server, it cannot see the newly created order, so it may return “not found.”

In the stated design, it is even worse than “before data is synced”: no shared storage or synchronization mechanism exists, so each server has its own isolated source of truth. The order may remain visible only through the server that received the write.

That is a distributed-systems correctness problem because the system must coordinate data across independently running machines to preserve the illusion of one order service.