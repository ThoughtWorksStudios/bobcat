create table Suppliers (
    id BIGSERIAL PRIMARY KEY,
    CompanyName     varchar(255),
    ContactName     varchar(255),
    ContactTitle    varchar(255),
    Address         varchar(255),
    City            varchar(255),
    Region          varchar(255),
    PostalCode      varchar(255),
    Country         varchar(255),
    Phone           varchar(255),
    Fax             varchar(255),
    HomePage        varchar(255)
);

create table Categories (
    id BIGSERIAL PRIMARY KEY,
    CategoryName    varchar(255),
    Description     varchar(255),
    Picture         integer
);

create table Products (
    id BIGSERIAL PRIMARY KEY,
    ProductName     varchar(255),
    SupplierID      INTEGER REFERENCES Suppliers (id),
    CategoryID      INTEGER REFERENCES Categories (id),
    QuantityPerUnit integer,
    UnitPrice       integer,
    UnitsInStock    integer,
    UnitsOnOrder    integer,
    Discontinued    boolean
);

create table Regions (
    id BIGSERIAL PRIMARY KEY,
    RegionDesc      varchar(255)
);

create table Territories (
    id BIGSERIAL PRIMARY KEY,
    TerritoryDesc   varchar(255),
    RegionID        INTEGER REFERENCES Regions (id)
);

create table Employees (
    id BIGSERIAL PRIMARY KEY,
    LastName        varchar(255),
    FirstName       varchar(255),
    Title           varchar(255),
    TitleOfCourtesy varchar(255),
    BirthDate       date,
    HireDate        date,
    Address         varchar(255),
    City            varchar(255),
    Region          varchar(255),
    PostalCode      varchar(255),
    Country         varchar(255),
    HomePhone       varchar(255),
    Extension       integer,
    Photo           integer,
    Notes           varchar(255),
    ReportsTo       INTEGER REFERENCES Employees (id),
    PhotoPath       varchar(255)
);

create table EmployeeTerritories (
    id BIGSERIAL PRIMARY KEY,
    EmployeeID      INTEGER REFERENCES Employees (id),
    TerritoryID     INTEGER REFERENCES Territories (id)
);

create table Customers (
    id BIGSERIAL PRIMARY KEY,
    CompanyName    varchar(255),
    ContactName    varchar(255),
    ContactTitle   varchar(255),
    Address        varchar(255),
    City           varchar(255),
    Region         varchar(255),
    PostalCode     varchar(255),
    Country        varchar(255),
    Phone          varchar(255),
    Fax            varchar(255)
);

create table Shippers (
    id BIGSERIAL PRIMARY KEY,
    CompanyName     varchar(255),
    Phone           varchar(255)
);

create table Orders (
    id BIGSERIAL PRIMARY KEY,
    CustomerID      INTEGER REFERENCES Customers (id),
    EmployeeID      INTEGER REFERENCES Employees (id),
    OrderDate       date,
    RequiredDate    date,
    ShippedDate     date,
    ShipVia         INTEGER REFERENCES Shippers (id),
    Freight         integer,
    ShipName        varchar(255),
    ShipAddress     varchar(255),
    ShipCity        varchar(255),
    ShipRegion      varchar(255),
    ShipPostalCode  varchar(255),
    ShipCountry     varchar(255)
);

create table OrderDetails (
    id BIGSERIAL PRIMARY KEY,
    OrderID         INTEGER REFERENCES Orders (id),
    ProductID       INTEGER REFERENCES Products (id),
    UnitPrice       integer,
    Quantity        integer,
    Discount        decimal
);

create table CustomerDemographics (
    id BIGSERIAL PRIMARY KEY,
    CustomerDesc   varchar(255)
);

create table CustomerCustomerDemo (
    id BIGSERIAL PRIMARY KEY,
    CustomerID      INTEGER REFERENCES Customers (id),
    CustomerDemographicsID  INTEGER REFERENCES CustomerDemographics (id)
);
