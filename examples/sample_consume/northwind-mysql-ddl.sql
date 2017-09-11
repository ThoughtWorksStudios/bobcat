DROP DATABASE IF EXISTS northwind;
CREATE DATABASE northwind;
USE northwind;

create table Suppliers (
    id INT NOT NULL PRIMARY KEY,
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
    id INT NOT NULL PRIMARY KEY,
    CategoryName    varchar(255),
    Description     varchar(255),
    Picture         integer
);

create table Products (
    id INT NOT NULL PRIMARY KEY,
    ProductName     varchar(255),
    SupplierID      INTEGER,
    FOREIGN KEY (SupplierID) REFERENCES Suppliers (id) ON DELETE CASCADE,
    CategoryID      INTEGER,
    FOREIGN KEY (CategoryID) REFERENCES Categories (id) ON DELETE CASCADE,
    QuantityPerUnit integer,
    UnitPrice       integer,
    UnitsInStock    integer,
    UnitsOnOrder    integer,
    Discontinued    boolean
);

create table Regions (
    id INT NOT NULL PRIMARY KEY,
    RegionDesc      varchar(255)
);

create table Territories (
    id INT NOT NULL PRIMARY KEY,
    TerritoryDesc   varchar(255),
    RegionID        INTEGER,
    FOREIGN KEY (RegionID) REFERENCES Regions (id) ON DELETE CASCADE
);

create table Employees (
    id INT NOT NULL PRIMARY KEY,
    LastName        varchar(255),
    FirstName       varchar(255),
    Title           varchar(255),
    TitleOfCourtesy varchar(255),
    BirthDate       datetime,
    HireDate        datetime,
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
    FOREIGN KEY (ReportsTo) REFERENCES Employees (id) ON DELETE CASCADE,
    PhotoPath       varchar(255)
);

create table EmployeeTerritories (
    id INT NOT NULL PRIMARY KEY,
    EmployeeID      INTEGER,
    FOREIGN KEY (EmployeeID) REFERENCES Employees (id) ON DELETE CASCADE,
    TerritoryID     INTEGER,
    FOREIGN KEY (TerritoryID) REFERENCES Territories (id) ON DELETE CASCADE
);

create table Customers (
    id INT NOT NULL PRIMARY KEY,
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
    id INT NOT NULL PRIMARY KEY,
    CompanyName     varchar(255),
    Phone           varchar(255)
);

create table Orders (
    id INT NOT NULL PRIMARY KEY,
    CustomerID      INTEGER,
    FOREIGN KEY (CustomerID) REFERENCES Customers (id) ON DELETE CASCADE,
    EmployeeID      INTEGER,
    FOREIGN KEY (EmployeeID) REFERENCES Employees (id) ON DELETE CASCADE,
    OrderDate       datetime,
    RequiredDate    datetime,
    ShippedDate     datetime,
    ShipVia         INTEGER,
    FOREIGN KEY (ShipVia) REFERENCES Shippers (id) ON DELETE CASCADE,
    Freight         integer,
    ShipName        varchar(255),
    ShipAddress     varchar(255),
    ShipCity        varchar(255),
    ShipRegion      varchar(255),
    ShipPostalCode  varchar(255),
    ShipCountry     varchar(255)
);

create table OrderDetails (
    id INT NOT NULL PRIMARY KEY,
    OrderID         INTEGER,
    FOREIGN KEY (OrderID) REFERENCES Orders (id) ON DELETE CASCADE,
    ProductID       INTEGER,
    FOREIGN KEY (ProductID) REFERENCES Products (id) ON DELETE CASCADE,
    UnitPrice       integer,
    Quantity        integer,
    Discount        decimal
);

create table CustomerDemographics (
    id INT NOT NULL PRIMARY KEY,
    CustomerDesc   varchar(255)
);

create table CustomerCustomerDemo (
    id INT NOT NULL PRIMARY KEY,
    CustomerID      INTEGER,
    FOREIGN KEY (CustomerID) REFERENCES Customers (id) ON DELETE CASCADE,
    CustomerDemographicsID  INTEGER,
    FOREIGN KEY (CustomerDemographicsID) REFERENCES CustomerDemographics (id) ON DELETE CASCADE
);
