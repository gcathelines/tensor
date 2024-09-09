
## Running the App
To run the app and the database, use the following command:
```shell
make run
```
By default GraphiQL is enabled, to disable it update **GRAPHIQL_ENABLED** env in `docker-compose.yaml`.

To run the migration for the target database, use the following command:
```shell
make migrate DSN=your://cool:dns@in:5432/here
```

## Testing the App
To run the app test, use the following command:
```shell
make test
```

## Project Structure

- `main.go`: The main file that orchestrates everything and runs the entire application.
- `config.yaml`: The configuration file used to run the app.
- `/config`: Contains the configuration structure needed for the app, mapping config from the `.yaml` file.
- `/graph`: Contains the generated files from GraphQL.
- `/internal`: Contains the main logic for the app.
- `/internal/database`: Contains the database logic, acting as a repository layer.
- `/internal/open_meteo`: Contains the client for the weather API by Open Meteo.
- `/internal/usecase`: Contains the main logic for the app, combining the database and weather API results to be presented in GraphQL.
- `/internal/types`: Contains the structs/models for objects in the API.
- `/migrations`: Contains the migration files.

## Planned Improvements

- Implement caching for power plant and weather API data to improve performance.
- Refactor the migration process to create a proper sequenced migration instead of using a Go script to run a SQL file directly on the target database.

## Additional Notes

- The weather API also returns elevation, but for this case, we will not be using the elevation data from the weather API. Instead, we will be using the elevation API as stated in the project requirements.
- The `hasPrecipitationToday` field is calculated using the daily precipitation sum. If the sum is greater than 0, it is marked as true.
- For simplicity, `BIGSERIAL` is chosen as the ID for power plants, as it provides a sortable ID for pagination. If dealing with a large amount of data and the possibility of running out of `int64` IDs, consider using [ULID](https://github.com/ulid/spec) instead. ULID is lexicographically sortable, ensuring correct pagination. Additionally, `LastID` is used instead of `offset` for pagination due to its better performance compared to offset-based pagination ([source](https://use-the-index-luke.com/sql/partial-results/fetch-next-page)).

