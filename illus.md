graph TD
    subgraph "Git Source"
        A[fa:fa-github GitHub Repository]
    end

    subgraph "TRIGRA Controller"
        B{fa:fa-bolt Webhook Handler}
        C[fa:fa-file-code Resource Parser]
        D[fa:fa-microchip Smart Applier]
    end

    subgraph "Infrastructure"
        E[fa:fa-dharmachakra K8s Cluster]
    end

    A -- "1. Push Event (Webhook)" --> B
    B -- "2. Fetch Manifests" --> A
    B --> C
    C --> D
    D -- "3. Atomic Apply" --> E

    %% Styling
    style B fill:#10b981,stroke:#333,stroke-width:2px,color:#fff
    style D fill:#10b981,stroke:#333,stroke-width:2px,color:#fff
    style A fill:#f1f5f9,stroke:#64748b
    style E fill:#f1f5f9,stroke:#64748b