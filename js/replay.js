// ===========================================
// REPLAY PAGE JAVASCRIPT
// =========================================== 

class GameReplay {
    constructor() {
        this.initializeElements();
        this.initializeBoard();
        this.setupEventListeners();
        
        // State tracking
        this.isLoading = false;
        this.boardToPlayIn = '?'; // Track which board should be highlighted
        this.moveCount = 0;
    }

    initializeElements() {
        // Board and controls
        this.boardEl = document.getElementById('board');
        this.prevBtn = document.getElementById('prev-btn');
        this.nextBtn = document.getElementById('next-btn');
        this.errorEl = document.getElementById('error-message');
        this.moveInfoEl = document.getElementById('move-info');
    }

    initializeBoard() {
        // Create 9 mini-boards with 81 cells (same structure as game board)
        this.boardEl.innerHTML = '';

        for (let mb = 0; mb < 9; mb++) {
            const miniBoard = document.createElement('div');
            miniBoard.className = 'mini-board';
            miniBoard.id = `mini-board-${mb}`;
            
            for (let cell = 0; cell < 9; cell++) {
                const globalIndex = mb * 9 + cell;
                const cellDiv = document.createElement('div');
                cellDiv.className = 'cell';
                cellDiv.id = `cell-${globalIndex}`;
                miniBoard.appendChild(cellDiv);
            }

            this.boardEl.appendChild(miniBoard);
        }
    }

    setupEventListeners() {
        this.prevBtn.addEventListener('click', () => this.previousMove());
        this.nextBtn.addEventListener('click', () => this.nextMove());
    }
    //
    // Get CSRF token from hidden field
    getCSRFToken() {
        const csrfInput = document.querySelector('#csrf-debug input[name="gorilla.csrf.Token"]') ||
                         document.querySelector('#csrf-debug input[name="_token"]') ||
                         document.querySelector('#csrf-debug input[type="hidden"]');

        if (csrfInput) {
            return csrfInput.value;
        }

        const csrfMeta = document.querySelector('meta[name="csrf-token"]');
        if (csrfMeta) {
            return csrfMeta.getAttribute('content');
        }

        return null;
    }

    // Show/hide loading state
    setLoading(loading) {
        this.isLoading = loading;
        this.prevBtn.disabled = loading;
        this.nextBtn.disabled = loading;

        if (loading) {
            this.moveInfoEl.innerHTML = '<span>Loading...</span>';
        }
    }

    // Show error message
    showError(message) {
        this.errorEl.textContent = `Error: ${message}`;
        this.errorEl.classList.remove('hidden');

        setTimeout(() => {
            this.errorEl.classList.add('hidden');
        }, 5000);
    }

    // Update move info display
    updateMoveInfo(moveDescription = '') {
        if (moveDescription) {
            this.moveInfoEl.innerHTML = `<span>${moveDescription}</span>`;
        } else {
            this.moveInfoEl.innerHTML = '<span>Use the controls above to navigate through the game</span>';
        }
    }

    // Send replay request to server
    async sendReplayRequest(move) {
        if (this.isLoading) return;

        console.log('Sending replay request with move:', move);
        this.setLoading(true);

        try {
            const csrfToken = this.getCSRFToken();
            const headers = {
                'Content-Type': 'application/json',
            };

            if (csrfToken) {
                headers['X-CSRF-Token'] = csrfToken;
            }

            const response = await fetch('/replay', {
                method: 'POST',
                headers: headers,
                body: JSON.stringify({ msg: move })
            });

            if (!response.ok) {
                let errorText = `HTTP ${response.status}`;
                try {
                    const errorData = await response.json();
                    errorText = errorData.message || errorText;
                } catch {
                    errorText = await response.text() || errorText;
                }
                throw new Error(errorText);
            }

            const data = await response.json();
            console.log('Response data:', data);

            // Update the board with response data
            if (data.type && Array.isArray(data.board)) {
                this.renderBoard(data.board);

                // Update mini-board winners if available
                if (data.complete_boards) {
                    this.updateMiniboardWinners(data.complete_boards);
                }

                // Update board highlighting if available
                if (data.board_to_play_in !== undefined) {
                    const boardByte = data.board_to_play_in;
                    if (boardByte === 63) { // '?' is ASCII 63
                        this.boardToPlayIn = '?';
                    } else if (boardByte >= 0 && boardByte <= 8) {
                        this.boardToPlayIn = boardByte.toString();
                    } else {
                        this.boardToPlayIn = '?';
                    }
                    this.highlightActiveBoard();
                } else {
                    console.log("data.board_to_play_in IS UNDEFINED")
                }

                // Update move counter and info
                if (move === 62) { // next move
                    this.moveCount++;
                    this.updateMoveInfo(`Move ${this.moveCount}`);
                } else if (move === 60) { // previous move
                    this.moveCount = Math.max(0, this.moveCount - 1);
                    this.updateMoveInfo(this.moveCount === 0 ? 'Start position' : `Move ${this.moveCount}`);
                }
            } else {
                throw new Error('Invalid response format');
            }

        } catch (error) {
            console.error('Replay request failed:', error);
            this.showError(error.message);
        } finally {
            this.setLoading(false);
        }
    }

    // Handle next move button
    nextMove() {
        this.sendReplayRequest(62); // ASCII for '>'
    }

    // Handle previous move button  
    previousMove() {
        this.sendReplayRequest(60); // ASCII for '<'
    }

    // Render board state (same as game board)
    renderBoard(state) {
        if (state.length !== 81) return;

        for (let i = 0; i < 81; i++) {
            const cell = document.getElementById(`cell-${i}`);
            const value = state[i] || '';

            cell.textContent = value.toUpperCase();

            // Add CSS classes for X and O
            cell.classList.remove('x', 'o');
            if (value.toLowerCase() === 'x') {
                cell.classList.add('x');
            } else if (value.toLowerCase() === 'o') {
                cell.classList.add('o');
            }
        }
    }

    // Update mini-board winners (same logic as game board)
    updateMiniboardWinners(completeBoardsBase64) {
        console.log('Updating mini-board winners:', completeBoardsBase64);

        let decodedBytes;
        try {
            const binaryString = atob(completeBoardsBase64);
            decodedBytes = new Uint8Array(binaryString.length);
            for (let i = 0; i < binaryString.length; i++) {
                decodedBytes[i] = binaryString.charCodeAt(i);
            }
        } catch (error) {
            console.error('Failed to decode complete_boards base64:', error);
            return;
        }

        for (let i = 0; i < 9; i++) {
            const miniBoard = document.getElementById(`mini-board-${i}`);
            const existingWinner = miniBoard.querySelector('.mini-board-winner');

            // Remove existing winner overlay
            if (existingWinner) {
                existingWinner.remove();
            }

            // Check if this mini-board is completed
            if (decodedBytes && i < decodedBytes.length) {
                const winnerByte = decodedBytes[i];

                let winner = null;
                if (winnerByte === 120) { // 'x' is ASCII 120
                    winner = 'x';
                } else if (winnerByte === 111) { // 'o' is ASCII 111
                    winner = 'o';
                } else if (winnerByte === 0) {
                    // Board not completed
                    miniBoard.classList.remove('completed');
                    continue;
                }

                if (winner) {
                    // Add completed class
                    miniBoard.classList.add('completed');

                    // Create winner overlay
                    const winnerOverlay = document.createElement('div');
                    winnerOverlay.className = `mini-board-winner ${winner}-winner`;
                    winnerOverlay.textContent = winner.toUpperCase();
                    miniBoard.appendChild(winnerOverlay);
                } else {
                    // Remove completed class
                    miniBoard.classList.remove('completed');
                }
            } else {
                // Remove completed class if no data
                miniBoard.classList.remove('completed');
            }
        }
    }

    // Highlight active board (same logic as game board)
    highlightActiveBoard() {
        console.log('Highlighting board. BoardToPlayIn:', this.boardToPlayIn);

        // Remove all existing highlights
        const allMiniBoards = document.querySelectorAll('.mini-board');
        allMiniBoards.forEach(board => board.classList.remove('highlight'));

        const ultimateBoard = document.querySelector('.ultimate-board');
        ultimateBoard.classList.remove('highlight-all');

        // Add highlights based on boardToPlayIn
        if (this.boardToPlayIn === '?') {
            // Highlight entire board
            ultimateBoard.classList.add('highlight-all');
        } else {
            // Parse board index and highlight specific mini-board
            const boardIndex = parseInt(this.boardToPlayIn);
            if (boardIndex >= 0 && boardIndex < 9) {
                const targetBoard = document.getElementById(`mini-board-${boardIndex}`);
                if (targetBoard) {
                    targetBoard.classList.add('highlight');
                }
            }
        }
    }
}

// Initialize the replay when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    new GameReplay();
});
