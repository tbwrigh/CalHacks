import sqlite3
import hashlib

# Connect to an SQLite database (in-memory for this example)
conn = sqlite3.connect(':memory:')
cursor = conn.cursor()

# Create a sample table
cursor.execute('''
    CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT,
        password TEXT
    )
''')

# Insert some dummy data
cursor.execute('INSERT INTO users (username, password) VALUES (?, ?)', ('admin', 'secret'))
cursor.execute('INSERT INTO users (username, password) VALUES (?, ?)', ('guest', 'guestpass'))
conn.commit()

def get_user_password(username):
    # This is vulnerable to SQL injection!
    query = f"SELECT password FROM users WHERE username = '{username}'"
    cursor.execute(query)
    result = cursor.fetchone()
    
    if result:
        return result[0]
    else:
        return None

# Vulnerable input - assume an attacker can input anything
user_input = input("Enter your username: ")

# Fetch and hash the password before printing (fixing the vulnerability)
password = get_user_password(user_input)
if password:
    hashed_password = hashlib.sha256(password.encode()).hexdigest()
    print(f"Hashed password for {user_input}: {hashed_password}")
else:
    print("User not found.")

# Close the connection
conn.close()