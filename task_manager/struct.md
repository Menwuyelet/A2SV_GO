task-manager/
├── Delivery/
│   ├── main.go
│   ├── controllers/
│   │   └── controller.go
│   └── routers/
│       └── router.go
├── Domain/
│   └── domain.go
├── Infrastructure/
│   ├── auth_middleWare.go
│   ├── jwt_service.go
│   └── password_service.go
├── Repositories/
│   ├── task_repository.go
│   └── user_repository.go
└── Usecases/
    ├── task_usecases.go
    └── user_usecases.go

Delivery/: Contains files related to the delivery layer, handling incoming requests and responses.
    main.go: Sets up the HTTP server, initializes dependencies, and defines the routing configuration.
   
    controllers/controllers.go: Handles incoming HTTP requests and invokes the appropriate use case methods.
    
    routers/routers.go: Sets up the routes and initializes the Gin router.

Domain/: Defines the core business entities and logic.
    
    domain.go: Contains the core business entities such as Task and User structs. 

Infrastructure/: Implements external dependencies and services.
    auth_middleWare.go: Middleware to handle authentication and authorization using JWT tokens.
    
    jwt_service.go: Functions to generate and validate JWT tokens.
    
    password_service.go: Functions for hashing and comparing passwords to ensure secure storage of user credentials.

Repositories/: Abstracts the data access logic.
    task_repository.go: Interface and implementation for task data access operations.
    
    user_repository.go: Interface and implementation for user data access operations.

Usecases/: Contains the application-specific business rules.
    task_usecases.go: Implements the use cases related to tasks, such as creating, updating, retrieving, and deleting tasks.
    
    user_usecases.go: Implements the use cases related to users, such as registering, logging in.
