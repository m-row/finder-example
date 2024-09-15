# Go Finder Example

this repo showcases the usage of [finder](https://github.com/m-row/finder)
for creating crud rest APIs with go & postgres.

---

## Commands

this project utilises docker to run `migration`

use `make help` to view command list:

---

## How to use

1. copy `.env.example` into a `.env` file and modify database connection info
2. run `make migrate/up/all` to create the `categories` table in the database
3. use [air](https://github.com/air-verse/air) to run the application
    - or use `go run .`

it features 5 actions inspired by laravel:

use base url: `localhost:8000/api/v1`

1. `index`: a list of the model, can be:
    - searched:
        - the columns used for search are inside `meta.search_columns`
        - `/categories?q=main`
        - search can utilize relational joins fields using the join alias
        when defined as:

        ```go
        func (m *Model) SearchFields() *[]string {
            return &[]string{
                "description",
                "rel:stores.name",
                "rel:stores.description",
            }
        }
        ```

    - paginated
        - pagination info are inside `meta` proprty of the response
        - `/categories?page=1&paginate=12`
    - sorted:
        - the columns used for sorting are inside `meta.columns`
        - ascending or descending order, and can use multiple column sort
        - `/categories?sorts=-created_at`
        - `/categories?sorts=created_at`
        - `/categories?sorts=depth,-created_at`
    - filtered:
        - the columns used for filtering are inside `meta.columns`
        - `/categories?filters=parent_id:null`
        - `/categories?filters=parent_id:not-null`
        - filter can use multiple `column` and multiple `criteria` values
        - `/categories?filters=parent_id:null,is_disabled:false|true`
        - filter can use `ops` between `column:op:criteria`, ops options below:
            - `eq`: db equal of `=` eg: `is_disabled:eq:true`
            - `nq`: db equal of `!=` eg: `is_disabled:nq:true`
            - `gt`: db equal of `>` eg: `price:gt:5`
            - `gte`: db equal of `>=` eg: `price:gte:5`
            - `lt`: db equal of `<` eg: `price:lt:5`
            - `lte`: db equal of `<=` eg: `price:lte:5`
            - `ex`: special filter for single relations of model eg: `roles:ex:5`
        - another eg: `/orders?filters=status:eq:completed|cancelled|returned`
    - mix:
        - all the above can be used in combination with each other
    - note:
        - `finder.Model` interface can use `Relations()` to define relations
        that can be used with `q` or `filters` for example if you have a `users`
        table with `roles` many-to-many relation:

        ```go
        func (m *Model) Relations() *[]finder.RelationField {
            return &[]finder.RelationField{
                {
                    Table: "roles",
                    Join: &finder.Join{
                        From: "users.id",
                        To:   "roles.id",
                    },
                    Through: &finder.Through{
                        Table: "user_roles",
                        Join: &finder.Join{
                            From: "user_roles.user_id",
                            To:   "user_roles.role_id",
                        },
                    },
                },
            }
        }
        ```

2. `show`: `GET /categories/:id` returns single model result
3. `destroy`: `DELETE /categories/:id` deletes single model
4. `store`: `POST /categories` creates single model
    - sample main category body:

    ```json
    {
        "id": "4db02628-50a2-4022-ac82-1949853bd728"
        "name": {
            "en":"main",
            "ar":"رئيسية"
        }
    }
    ```

    - sample sub category body:

    ```json
    {
        "name": {
            "en":"branch",
            "ar":"فرعية"
        },
        "is_disabled":true,
        "is_featured":true,
        "parent":{
            "id": "4db02628-50a2-4022-ac82-1949853bd728"
        }
    }
    ```

5. `update`: `PUT /categories/:id` updates single model result
    - sample update category body:

    ```json
    {
        "name": {
            "en":"branch name"
        },
        "is_disabled":false,
        "is_featured":true
    }
    ```

---

## resources

- [pgx](https://github.com/jackc/pgx) Replaces the standard `database/sql` driver
- [sqlx](https://github.com/jmoiron/sqlx) Extends the standard `database/sql` functions
- [echo](https://github.com/labstack/echo) minimalist web framework
- [squirrel](https://github.com/Masterminds/squirrel) Query Builder
