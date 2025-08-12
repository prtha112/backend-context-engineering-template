## FEATURE:

Product CRUD (Create, Read, Update, Delete)
```bash
# Table structure
product {
    id int serial primary key
    store_id integer <-- please not add foreign key to this. just keep normal fix integer 
    name varchar(100)
    description text
    amount integer
    price numeric(12,2)
}
```

## EXAMPLES:

In the examples/ folder, there is a README for you to read to understand what the example is all about and also how to structure your own README when you create documentation for the above feature.

## DOCUMENTATION:

Documentation: https://go.dev/doc/

## OTHER CONSIDERATIONS:

-