package main;

// == IMPORTS

import ( "database/sql"
    "fmt"
    "strconv"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gocql/gocql"
    "io/ioutil"
    "os"
    "strings"
    "time" );

// -- VARIABLES

var IsCqlDatabase, IsSqlDatabase bool;
var DatabaseDriver, DatabasePassword, DatabasePort, DatabaseServer, DatabaseSchema, DatabaseUser string;
var ScriptFilePathArray [] string;
var CqlSession * gocql.Session;
var SqlDatabase * sql.DB;

// -- TYPES

type ERROR_MESSAGE struct {
    Text string;
    Error error;
}

// .. INQUIRIES

func ( error_message * ERROR_MESSAGE ) Print( ) {
    var error_text string;

    if ( error_message != nil ) {
        text := error_message.Text;

        if ( error_message.Error != nil ) {
            error_text = error_message.Error.Error();
        }

        if ( text != "" ) {
             if ( error_text != "" ) {
                fmt.Println( text, "(" + error_text + ")" );
            } else {
                fmt.Println( text );
            }
        } else if ( error_text != "" ) {
            fmt.Println( error_text );
        }
    }
}

// .. OPERATIONS

func ( error_message * ERROR_MESSAGE ) Set( text string, error_ error ) {
    error_message.Text = text;
    error_message.Error = error_;
}

// ~~

func ( error_message * ERROR_MESSAGE ) SetText( text string ) {
    error_message.Text = text;
    error_message.Error = nil;
}

// ~~

func ( error_message * ERROR_MESSAGE ) SetError( error_ error ) {
    error_message.Text = "";
    error_message.Error = error_;
}

// -- FUNCTIONS

func IsNatural( text string ) bool {
    if ( len( text ) == 0 ) {
        return false;
    } else {
        for _, character := range text {
            if ( character < '0' || character > '9' ) {
                return false;
            }
        }

        return true;
    }
}

// ~~

func GetInteger( text string ) int {
    integer, _ := strconv.ParseInt( text, 10, 64 );

    return int( integer );
}

// ~~

func OpenDatabase( error_message * ERROR_MESSAGE ) bool {
    var error_ error;

    fmt.Println( "Opening database." );

    if ( IsCqlDatabase ) {
        cql_cluster_configuration := gocql.NewCluster( DatabaseServer );
        cql_cluster_configuration.Keyspace = DatabaseSchema;
        cql_cluster_configuration.Port = GetInteger( DatabasePort );
        cql_cluster_configuration.Timeout = 15 * time.Second;
        cql_cluster_configuration.ConnectTimeout = 15 * time.Second;
        cql_cluster_configuration.Consistency = gocql.Quorum;

        CqlSession, error_ = cql_cluster_configuration.CreateSession();
    } else if ( IsSqlDatabase ) {
        SqlDatabase, error_ = sql.Open( "mysql", DatabaseUser + ":" + DatabasePassword + "@tcp(" + DatabaseServer + ":" + DatabasePort + ")/" );
    }

    if ( error_ != nil ) {
        error_message.SetError( error_ );

        return false;
    } else {
        return true;
    }
}

// ~~

func RunDatabaseQuery( query string, error_message * ERROR_MESSAGE ) bool {
    var error_ error;

    fmt.Println( query );

    if ( IsCqlDatabase ) {
        error_ = CqlSession.Query( query ).Exec();
    } else if ( IsSqlDatabase ) {
        _, error_ = SqlDatabase.Exec( query );
    }

    if ( error_ != nil ) {
        error_message.SetError( error_ );

        return false;
    } else {
        return true;
    }
}

// ~~

func ExecuteScripts( error_message * ERROR_MESSAGE ) bool {
    for _, script_file_path := range ScriptFilePathArray {
        fmt.Println( "Reading file : " + script_file_path )

        byte_array, error_ := ioutil.ReadFile( script_file_path );

        if ( error_ != nil ) {
            error_message.SetError( error_ );

            return false;
        }

        script := strings.ReplaceAll( string( byte_array ), "\r", "" );
        line_array := strings.Split( script, "\n" );
        query := "";

        for _, line := range line_array {
            if ( len( line ) > 0 ) {
                query += line;

                if ( strings.HasSuffix( query, ";" ) ) {
                    if ( !RunDatabaseQuery( query, error_message ) ) {
                        return false;
                    }

                    query = "";
                }
            }
        }
    }

    return true;
}

// ~~

func CloseDatabase( ) bool {
    if ( IsCqlDatabase ) {
        CqlSession.Close();
    } else if ( IsSqlDatabase ) {
        SqlDatabase.Close();
    }

    return true;
}

// ~~

func ParseArguments( error_message * ERROR_MESSAGE ) bool {
    argument_array := os.Args[ 1 : ];

    if ( len( argument_array ) >= 7 ) {
        DatabaseDriver = argument_array[ 0 ];
        DatabaseServer = argument_array[ 1 ];
        DatabasePort = argument_array[ 2 ];
        DatabaseSchema = argument_array[ 3 ];
        DatabaseUser = argument_array[ 4 ];
        DatabasePassword = argument_array[ 5 ];
        ScriptFilePathArray = argument_array[ 6 : ];

        fmt.Println( "Driver :", DatabaseDriver );
        fmt.Println( "Server :", DatabaseServer );
        fmt.Println( "Port :", DatabasePort );
        fmt.Println( "Schema :", DatabaseSchema );
        fmt.Println( "User :", DatabaseUser );
        fmt.Println( "Password :", DatabasePassword );

        if ( DatabaseDriver == "cassandra" ) {
            IsCqlDatabase = true;
        } else if ( DatabaseDriver == "mysql" ) {
            IsSqlDatabase = true;
        } else {
            error_message.SetText( "Invalid database driver : " + DatabaseDriver );

            return false;
        }

        if ( DatabaseServer == "" ) {
            error_message.SetText( "Invalid database server : " + DatabaseServer );

            return false;
        }

        if ( DatabasePort == "" || !IsNatural( DatabasePort ) ) {
            error_message.SetText( "Invalid database port : " + DatabasePort );

            return false;
        }

        if ( DatabaseSchema == "" ) {
            error_message.SetText( "Missing database name argument : " + DatabaseSchema );

            return false;
        }

        if ( DatabaseUser == "" ) {
            error_message.SetText( "Missing database name argument : " + DatabaseUser );

            return false;
        }


        for _, script_file_path := range ScriptFilePathArray {
            if ( ( strings.HasSuffix( script_file_path, ".cql" ) && IsCqlDatabase ) || ( strings.HasSuffix( script_file_path, ".sql" ) && IsSqlDatabase ) ) {
                fmt.Println( "Script :", script_file_path );
            } else {
                error_message.SetText( "Invalid script argument : " + script_file_path );

                return false;
            }
        }
    } else {
        error_message.SetText( "Missing arguments" );

        return false;
    }

    return true;
}

// ~~

func main( ) {
    var error_message ERROR_MESSAGE;

    if ( ParseArguments( &error_message ) && OpenDatabase( &error_message ) && ExecuteScripts( &error_message ) && CloseDatabase() ) {
        fmt.Println( "Success." );
    } else {
        error_message.Print();
    }
}
