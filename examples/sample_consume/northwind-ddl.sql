create table Suppliers (
id INTEGER NOT NULL,
    SupplierID      integer,
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
    HomePage        varchar(255),
    PRIMARY KEY (ID)
);


create table Categories (
id INTEGER NOT NULL,
    CategoryID      integer,
    CategoryName    varchar(255),
    Description     varchar(255),
    Picture         integer,
    PRIMARY KEY (ID)
);



create table Products (
id INTEGER NOT NULL,
    ProductID       integer,
    ProductName     varchar(255),
    SupplierID      INTEGER,
    CategoryID      INTEGER,
    QuantityPerUnit integer,
    UnitPrice       integer,
    UnitsInStock    integer,
    UnitsOnOrder    integer,
    Discontinued    boolean,
    PRIMARY KEY (ID)
);



create table Region (
id INTEGER NOT NULL,
    RegionID        integer,
    RegionDesc      varchar(255),
    PRIMARY KEY (ID)
);



create table Territories (
id INTEGER NOT NULL,
    TerritoryID     integer,
    TerritoryDesc   varchar(255),
    RegionID        INTEGER,
    PRIMARY KEY (ID)
);



create table Employees (
id INTEGER NOT NULL,
    EmployeeID      integer,
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
    ReportsTo       INTEGER,
    PhotoPath       varchar(255),
    PRIMARY KEY (ID)
);




create table EmployeeTerritories (
id INTEGER NOT NULL,
    EmployeeID      INTEGER,
    TerritoryID     INTEGER,
    PRIMARY KEY (ID)
);

create table Customers (
id INTEGER NOT NULL,
    CustomerID     integer,
    CompanyName    varchar(255),
    ContactName    varchar(255),
    ContactTitle   varchar(255),
    Address        varchar(255),
    City           varchar(255),
    Region         varchar(255),
    PostalCode     varchar(255),
    Country        varchar(255),
    Phone          varchar(255),
    Fax            varchar(255),
    PRIMARY KEY (ID)
);



create table Shippers (
id INTEGER NOT NULL,
    ShipperID       integer,
    CompanyName     varchar(255),
    Phone           varchar(255),
    PRIMARY KEY (ID)
);



create table Orders (
id INTEGER NOT NULL,
    OrderID         integer,
    CustomerID      INTEGER,
    EmployeeID      INTEGER,
    OrderDate       date,
    RequiredDate    date,
    ShippedDate     date,
    ShipVia         INTEGER,
    Freight         integer,
    ShipName        varchar(255),
    ShipAddress     varchar(255),
    ShipCity        varchar(255),
    ShipRegion      varchar(255),
    ShipPostalCode  varchar(255),
    ShipCountry     varchar(255),
    PRIMARY KEY (ID)
);


create table OrderDetails (
id INTEGER NOT NULL,
    OrderID         INTEGER,
    ProductID       INTEGER,
    UnitPrice       integer,
    Quantity        integer,
    Discount        decimal,
    PRIMARY KEY (ID)
);

create table CustomerDemographics (
id INTEGER NOT NULL,
    CustomerTypeID integer,
    CustomerDesc   varchar(255),
    PRIMARY KEY (ID)
);



create table CustomerCustomerDemo (
id INTEGER NOT NULL,
    CustomerID      INTEGER,
    CustomerTypeID  INTEGER,
    PRIMARY KEY (ID)
);




