// ===========================================
// DELETE ACCOUNT PAGE JAVASCRIPT
// =========================================== 

document.addEventListener('DOMContentLoaded', function() {
    const userId = '{{.user_id}}';
    const deleteBtn = document.getElementById('deleteBtn');
    const deleteText = document.getElementById('deleteText');
    const loadingSpinner = document.getElementById('loadingSpinner');

    deleteBtn.addEventListener('click', async (e) => {
        e.preventDefault();

        // Add loading state
        deleteBtn.disabled = true;
        deleteBtn.classList.add('loading');

        // Extract CSRF token if present
        const csrfInput = document.querySelector('#csrf-debug input[name="csrf_token"]');
        const csrfToken = csrfInput ? csrfInput.value : null;

        try {
            console.log(`Sending DELETE request to: /profile/${userId}`);
            const res = await fetch(`/profile/${userId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                    ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {})
                }
            });

            console.log('Response status:', res.status);

            if (res.status === 200) {
                console.log('Account deleted - redirecting to home');
                window.location.href = '/';
            } else {
                console.warn('Delete failed - redirecting to error page');
                window.location.href = '/error-page';
            }
        } catch (err) {
            console.error('Fetch error:', err);
            alert('Request failed: ' + err.message);
            
            // Remove loading state on error
            deleteBtn.disabled = false;
            deleteBtn.classList.remove('loading');
        }
    });
});
