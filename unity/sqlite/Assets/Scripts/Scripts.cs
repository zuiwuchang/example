using UnityEngine;
using UnityEngine.UI;
using Mono.Data.Sqlite;
public static class DbConnectionExtend
{
    public delegate T Execute<T>();
    public delegate void Execute();
    public static T Transaction<T>(this SqliteConnection db, Execute<T> f)
    {
        SqliteTransaction transaction = db.BeginTransaction();
        try
        {
            T result = f();
            transaction.Commit();
            return result;
        }
        catch (System.Exception)
        {
            transaction.Rollback();
            throw;
        }
    }
    public static void Transaction(this SqliteConnection db, Execute f)
    {
        SqliteTransaction transaction = db.BeginTransaction();
        try
        {
            f();
            transaction.Commit();
        }
        catch (System.Exception)
        {
            transaction.Rollback();
            throw;
        }
    }
}


public class Scripts : MonoBehaviour
{
    public string value;
    public Text message;
    public string dbname = "test.db";
    public string connectionString
    {
        get
        {
            var str = System.IO.Path.Combine(Application.persistentDataPath, dbname);
            // str = ":memory:";
            str = "URI=file:" + str;
            return str;
        }
    }
    private SqliteConnection _conn;
    private void OnDestroy()
    {
        _conn?.Close();
    }
    private void SetError(System.Exception e)
    {
        var msg = $"error: {e}";
        Debug.LogWarning(msg);
        message.text = msg;
    }

    public void OnClickOpen()
    {
        if (_conn != null)
        {
            message.text = "error: already opened";
            return;
        }
        try
        {
            message.text = $"connecting: {connectionString}";
            _conn = new SqliteConnection(connectionString);
            _conn.Open();
            message.text = $"connected: {connectionString}";
        }
        catch (System.Exception e)
        {
            if (_conn != null)
            {
                _conn.Close();
                _conn = null;
            }
            SetError(e);
        }
    }
    public void OnClickClose()
    {
        if (_conn == null)
        {
            message.text = "error: already closed";
        }
        else
        {
            _conn.Close();
            _conn = null;
            message.text = "closed";
        }
    }
    public void OnClickCreateTable()
    {
        try
        {
            SqliteCommand commnad = _conn.CreateCommand();
            commnad.CommandText = "CREATE TABLE IF NOT EXISTS List (id INTEGER PRIMARY KEY AUTOINCREMENT,name VARCHAR(20) UNIQUE NOT NULL)";
            int rows = commnad.ExecuteNonQuery();
            message.text = $"Create List; rows={rows}";
        }
        catch (System.Exception e)
        {
            SetError(e);
        }
    }
    public void OnClickDropTable()
    {
        try
        {
            SqliteCommand commnad = _conn.CreateCommand();
            commnad.CommandText = "DROP TABLE IF EXISTS List";
            int rows = commnad.ExecuteNonQuery();
            message.text = $"DROP List; rows={rows}";
        }
        catch (System.Exception e)
        {
            SetError(e);
        }
    }
    public void OnClickInsert()
    {
        try
        {
            var id = _conn.Transaction(() =>
              {
                  SqliteCommand commnad = _conn.CreateCommand();
                  commnad.Parameters.AddWithValue("@name", value);
                  commnad.CommandText = "INSERT INTO List (name) VALUES (@name); SELECT LAST_INSERT_ROWID();";
                  var id = (long)commnad.ExecuteScalar();
                  return id;
              });
            message.text = $"INSERT id: {id}";
        }
        catch (System.Exception e)
        {
            SetError(e);
        }
    }
    public void OnClickQuery()
    {
        try
        {
            SqliteCommand commnad = _conn.CreateCommand();
            commnad.Parameters.AddWithValue("@name", value);
            commnad.CommandText = "SELECT id FROM List WHERE name = @name";
            SqliteDataReader reader = commnad.ExecuteReader();
            while (reader.Read())
            {
                // var id = reader["id"];
                var id = reader.GetInt64(0);
                message.text = $"{value}'s id is {id}";
                return;
            }
            message.text = $"not found: {value}";
        }
        catch (System.Exception e)
        {
            SetError(e);
        }
    }
    public void OnClickDelete()
    {
        try
        {
            SqliteCommand commnad = _conn.CreateCommand();
            commnad.CommandText = "DELETE FROM List WHERE name = @name";
            commnad.Parameters.AddWithValue("@name", value);
            int rows = commnad.ExecuteNonQuery();
            message.text = $"DELETE rows: {rows}";
        }
        catch (System.Exception e)
        {
            SetError(e);
        }
    }
}
