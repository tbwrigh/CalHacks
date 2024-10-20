```python
def get_user_password(username):
    query = "SELECT password FROM users WHERE username = ?"
    cursor.execute(query, (username,))
    result = cursor.fetchone()
    
    if result:
        return result[0]
    else:
        return None
```