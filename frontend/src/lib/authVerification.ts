const usernameRegex = /^[a-zA-Z0-9_.-]+$/
const passwordRegex = /^[a-zA-Z0-9!"#$%&'()*+,\-./:;<=>?@[\\\]^_{|}~]+$/

export const checkUsername = (username: string): string  => {
    if (username.length < 3) return "username must be 3 or more character long"
    if (username.length > 30) return "username must be 30 characters or less"
    if (!usernameRegex.test(username)) return "invalid character detected in username (can only container letters, numbers and _, . and -"
    return ""
}

export const checkPassword = (password: string): string => {
    if (password.length < 8) return "password must be 8 or more characters long"
    if (password.length > 64) return "password must be 64 characters or less"
    if (!passwordRegex.test(password)) return "invalid character detected in password"
    return ""
}