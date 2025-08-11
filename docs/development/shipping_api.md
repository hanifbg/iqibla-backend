# Shipping API Documentation

## Endpoints

### Get Provinces

Get a list of all provinces or a specific province by ID.

- **URL**: `/api/v1/shipping/provinces`
- **Method**: `GET`
- **Query Parameters**:
  - `id` (optional): Province ID to get a specific province
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Provinces retrieved successfully",
      "data": [
        {
          "province_id": "1",
          "province": "Bali"
        },
        ...
      ]
    }
    ```
- **Error Response**:
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to get provinces",
      "message": "Error details"
    }
    ```

### Get Cities

Get cities by province ID and optionally filter by city ID.

- **URL**: `/api/v1/shipping/cities/:province_id`
- **Method**: `GET`
- **URL Parameters**:
  - `province_id`: ID of the province
- **Query Parameters**:
  - `id` (optional): City ID to get a specific city
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Cities retrieved successfully",
      "data": [
        {
          "city_id": "1",
          "province_id": "1",
          "city_name": "Badung"
        },
        ...
      ]
    }
    ```
- **Error Response**:
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to get cities",
      "message": "Error details"
    }
    ```

### Get Districts

Get districts by city ID.

- **URL**: `/api/v1/shipping/districts/:city_id`
- **Method**: `GET`
- **URL Parameters**:
  - `city_id`: ID of the city
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Districts retrieved successfully",
      "data": [
        {
          "district_id": "1",
          "city_id": "575",
          "district_name": "Cengkareng"
        },
        ...
      ]
    }
    ```
- **Error Response**:
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to get districts",
      "message": "Error details"
    }
    ```

### Calculate Shipping Cost

Calculate shipping cost based on origin, destination, weight, and courier.

- **URL**: `/api/v1/shipping/cost`
- **Method**: `POST`
- **Request Body**:
  ```json
  {
    "origin": "501",
    "destination": "114",
    "weight": 1000,
    "courier": "jne"
  }
  ```
- **Success Response**:
  - **Code**: 200
  - **Content**:
    ```json
    {
      "message": "Shipping cost calculated successfully",
      "data": [
        {
          "service": "OKE",
          "description": "Ongkos Kirim Ekonomis",
          "cost": 38000,
          "etd": "2-3"
        },
        ...
      ]
    }
    ```
- **Error Response**:
  - **Code**: 400
  - **Content**:
    ```json
    {
      "error": "Invalid request",
      "message": "Error details"
    }
    ```
  - **Code**: 500
  - **Content**:
    ```json
    {
      "error": "Failed to calculate shipping cost",
      "message": "Error details"
    }
    ```
