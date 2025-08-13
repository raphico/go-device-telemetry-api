Modern IoT and robotics systems require end-to-end solutions: from embedded firmware collecting sensor data to backend services processing, storing, and responding to that data in real time.

This project was built to explore how to design and implement such a backend service in Go, focusing on:

- Device management: securely registering and authenticating hardware.
- Telemetry ingestion: collecting high-volume, time-series data for analysis and monitoring.
- Command dispatch: enabling remote control and configuration of devices.
- Scalable architecture: following clean, testable design patterns.

While inspired by real-world industrial IoT platforms, this implementation is intentionally lightweight so it can be deployed quickly for small-scale robotics, environmental monitoring, or smart home setups.

The goal is to demonstrate practical backend engineering skills—covering authentication, data modeling, CRUD APIs, automated testing, and CI/CD pipelines—while building something that could directly integrate with microcontroller firmware, robotics systems, or sensor networks.
