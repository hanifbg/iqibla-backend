# Ticket SRP-001: Create Repository Interface and Structure

## Summary
Create the foundational repository interface and directory structure for the shipping repository refactor.

## Description
Establish the repository layer foundation by creating the `ShippingRepository` interface and setting up the proper directory structure for the RajaOngkir implementation. Since the service is not yet released, we can implement this directly without backward compatibility concerns.

## Acceptance Criteria

### Must Have
- [x] Create `internal/repository/shipping.go` interface with all required methods
- [x] Create `internal/repository/rajaongkir/` directory structure
- [x] Define configuration structure for repository initialization
- [x] Set up mock generation directives
- [x] Ensure interface follows existing repository patterns in the codebase

### Should Have
- [x] Add comprehensive interface documentation
- [x] Include error type definitions
- [x] Set up configuration patterns for extensibility
- [x] Prepare for immediate integration (no compatibility layer needed)

### Could Have
- [x] Create foundation for metrics and observability
- [x] Prepare structure for caching layer integration

## Technical Requirements

### Interface Definition
```go
// internal/repository/shipping.go
package repository

import "github.com/hanifbg/landing_backend/internal/model/response"

//go:generate mockgen -source=shipping.go -destination=../service/shipping/mocks/shipping_repository_mock.go -package=mocks

type ShippingRepository interface {
    GetProvinces(provinceID string) ([]response.RajaOngkirProvince, error)
    GetCities(provinceID, cityID string) ([]response.RajaOngkirCity, error)
    GetDistricts(cityID string) ([]response.RajaOngkirDistrict, error)
    CalculateShippingCost(origin, destination string, weight int, courier string) ([]response.RajaOngkirCost, error)
}
```

### Directory Structure
```
internal/repository/
├── shipping.go (new interface)
└── rajaongkir/ (new directory)
    ├── client.go (moved and refactored)
    ├── client_test.go (comprehensive tests)
    ├── config.go (configuration types)
    └── options.go (optional patterns)
```

### Error Handling
```go
// Define repository-specific error types
type RepositoryError struct {
    Operation string
    Cause     error
    Details   map[string]interface{}
}
```

## Implementation Steps

1. **Create interface file** (`internal/repository/shipping.go`)
   - Define `ShippingRepository` interface
   - Add proper package documentation
   - Include mock generation directive

2. **Create directory structure**
   - Create `internal/repository/rajaongkir/` directory
   - Set up basic file structure

3. **Define configuration types** (`internal/repository/rajaongkir/config.go`)
   - Create configuration struct
   - Add validation functions
   - Define default values

4. **Set up optional patterns** (`internal/repository/rajaongkir/options.go`)
   - Create option functions for configuration
   - Prepare for future extensibility

5. **Add error types** 
   - Define repository-specific errors
   - Implement proper error wrapping

## Dependencies
- None (foundational ticket)

## Estimated Effort
- **Development**: 4 hours
- **Testing**: 2 hours
- **Documentation**: 1 hour
- **Total**: 7 hours

## Definition of Done
- [x] Repository interface compiles without errors
- [x] Directory structure is created and organized
- [x] Configuration types are properly defined
- [x] Mock generation works correctly
- [x] Code follows project style guidelines
- [x] Documentation is complete and accurate
- [x] All acceptance criteria are met

## Risk Assessment
- **Low Risk**: Foundational changes with no immediate impact on existing functionality
- **Mitigation**: Create alongside existing code without removing anything

## Notes
- This ticket creates the foundation for all subsequent refactoring work
- Interface should be designed to support future caching and multiple provider implementations
- Follow existing repository patterns established in the codebase
- Ensure consistency with other repository interfaces

## Security Considerations

### Interface Design
- [x] Design methods to require explicit parameters (no map[string]interface{})
- [x] Ensure interface methods encourage validation before external calls
- [x] Define error types that avoid leaking implementation details

### Sensitive Information Handling
- [x] Design configuration structures to support secure API key storage
- [x] Ensure error types don't expose credentials or secrets
- [x] Support structured logging with redaction of sensitive values

### Input Validation
- [x] Design methods to encourage proper parameter validation
- [x] Document validation requirements for interface implementation
- [x] Avoid interfaces that bypass input sanitization

## Related Tickets
- **Blocks**: SRP-002 (Move RajaOngkir Implementation)
- **Blocks**: SRP-003 (Update Service Layer)

## Labels
- `refactor`
- `architecture`
- `repository`
- `foundation`
