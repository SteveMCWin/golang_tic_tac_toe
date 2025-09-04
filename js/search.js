// ===========================================
// SEARCH PAGE JAVASCRIPT
// =========================================== 

async function findUsers() {
    const query = document.getElementById("searchQuery").value.trim();
    const resultsContainer = document.getElementById("results");
    resultsContainer.innerHTML = "";  // clear old results

    if (!query) {
        return;
    }

    try {
        const res = await fetch(`/search?name=` + encodeURIComponent(query));
        if (!res.ok) {
            resultsContainer.innerHTML = "<div class='message error'>Error searching users.</div>";
            return;
        }

        const users = await res.json();

        if (users.length === 0) {
            resultsContainer.innerHTML = "<div class='message no-results'>No users found.</div>";
            return;
        }

        for (const user of users) {
            const link = document.createElement("a");
            link.href = `/profile/${user.id}`;
            link.className = "user-result";
            
            // Create avatar image (only if avatar URL exists)
            if (user.avatar_url || user.avatarURL) {
                const avatar = document.createElement("img");
                avatar.src = user.avatar_url || user.avatarURL;
                avatar.alt = `${user.username} Avatar`;
                avatar.className = "user-avatar";
                link.appendChild(avatar);
            }
            
            // Create username span
            const username = document.createElement("span");
            username.textContent = user.username;
            username.className = "user-name";
            
            // Append username to link
            link.appendChild(username);
            
            resultsContainer.appendChild(link);
        }
    } catch (error) {
        console.error('Search error:', error);
        resultsContainer.innerHTML = "<div class='message error'>Error searching users.</div>";
    }
}

// Add Enter key support for search input
document.addEventListener('DOMContentLoaded', function() {
    const searchInput = document.getElementById('searchQuery');
    
    searchInput.addEventListener('keydown', function(e) {
        if (e.key === 'Enter') {
            findUsers();
        }
    });
});
