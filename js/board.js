// ===========================================
// GAME BOARD JAVASCRIPT
// =========================================== 

class GameBoard {
    constructor() {
        this.initializeElements();
        this.initializeWebSocket();
        this.initializeBoard();
        
        // Timer state
        this.p1Time = 0;
        this.p2Time = 0;
        this.lastUpdateTime = Date.now();
        this.currentPlayer = 1; // 1 for player1, 2 for player2
        this.timerInterval = null;
        this.timersStarted = false; // Track if timers have been started
        
        // Board state
        this.boardToPlayIn = '?'; // Start with all boards highlighted
        this.highlightActiveBoard();
    }
    
    initializeElements() {
        // Main containers
        this.boardEl = document.getElementById('board');
        this.waitingMessageEl = document.getElementById('waiting-message');
        this.gameContainerEl = document.getElementById('game-container');
        
        // Game header elements
        this.currentTurnEl = document.getElementById('current-turn');
        this.p1TimerEl = document.getElementById('p1-timer');
        this.p2TimerEl = document.getElementById('p2-timer');
        this.player1InfoEl = document.getElementById('player1-info');
        this.player2InfoEl = document.getElementById('player2-info');
        
        // Game over elements
        this.gameOverOverlayEl = document.getElementById('game-over-overlay');
        this.gameOverTitleEl = document.getElementById('game-over-title');
    }
    
    initializeWebSocket() {
        // Get game mode from the global config
        const gameMode = window.appConfig?.gameMode;
        
        if (gameMode === undefined || gameMode === null) {
            console.error('Game mode not found in appConfig');
            return;
        }
        
        console.log('Initializing WebSocket with game_mode:', gameMode);
        
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/ws?game_mode=${gameMode}`;
        
        console.log('Connecting to WebSocket:', wsUrl);
        
        this.socket = new WebSocket(wsUrl);
        
        this.socket.onopen = () => {
            console.log('WebSocket connection opened');
        };
        
        this.socket.onmessage = (event) => {
            console.log('WebSocket message received:', event.data);
            this.handleWebSocketMessage(event);
        };
        
        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
        
        this.socket.onclose = (event) => {
            console.log('WebSocket connection closed', event);
            this.stopTimer();
        };
    }
    
    initializeBoard() {
        // Create 9 mini-boards with 81 clickable cells
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
                cellDiv.addEventListener('click', () => this.handleCellClick(mb, cell));
                miniBoard.appendChild(cellDiv);
            }
            
            this.boardEl.appendChild(miniBoard);
        }
    }
    
    handleWebSocketMessage(event) {
        const msg = JSON.parse(event.data);
        
        switch (msg.type) {
            case 'found':
                this.handleOpponentFound(msg);
                break;
            case 'start':
                this.handleGameStart(msg);
                break;
            case 'state':
                this.handleGameState(msg);
                break;
            case 'finish':
                this.handleGameFinish(msg);
                break;
        }
    }
    
    handleOpponentFound(msg) {
        // Hide waiting message and show game
        this.waitingMessageEl.classList.add('hidden');
        this.gameContainerEl.classList.remove('hidden');
        
        // Update player names
        if (msg.p1name) {
            this.player1InfoEl.querySelector('.player-name').textContent = msg.p1name;
        }
        if (msg.p2name) {
            this.player2InfoEl.querySelector('.player-name').textContent = msg.p2name;
        }
        
        // Set timer display to --:-- and don't start timers yet
        this.p1TimerEl.textContent = '--:--';
        this.p2TimerEl.textContent = '--:--';
        this.timersStarted = false;
        
        // Set initial turn display (assuming player 1/X starts)
        this.currentTurnEl.textContent = "Player X's Turn";
        this.player1InfoEl.className = "player-info active";
        this.player2InfoEl.className = "player-info inactive";
        
        console.log('Opponent found:', msg);
    }
    
    handleGameStart(msg) {
        // Set the flag first, then update timers with server values and start them
        this.timersStarted = true;
        this.updateTimersFromServer(msg.p1_time, msg.p2_time, true);
        
        console.log('Game started:', msg);
    }
    
    handleGameState(msg) {
        if (Array.isArray(msg.board)) {
            this.renderBoard(msg.board);
        }
        
        // Only update timers if they have been started
        if (this.timersStarted) {
            // Update timers with server correction
            this.updateTimersFromServer(msg.p1_time, msg.p2_time, msg.p1_move);
        }
        
        // Update UI state
        this.updateGameUI(msg.p1_move);
        
        // Update mini-board winners
        if (msg.complete_boards) {
            this.updateMiniboardWinners(msg.complete_boards);
        }
        
        // Update board highlighting - handle byte value correctly
        if (msg.board_to_play_in !== undefined) {
            const boardByte = msg.board_to_play_in;
            console.log('Raw board_to_play_in byte value:', boardByte);
            
            if (boardByte === 63) { // '?' is ASCII 63
                this.boardToPlayIn = '?';
            } else if (boardByte >= 0 && boardByte <= 8) {
                this.boardToPlayIn = boardByte.toString();
            } else {
                console.warn('Unexpected board_to_play_in value:', boardByte);
                this.boardToPlayIn = '?';
            }
            
            console.log('Processed boardToPlayIn:', this.boardToPlayIn);
            this.highlightActiveBoard();
        }
    }
    
    handleGameFinish(msg) {
        this.stopTimer();
        
        let title = '';
        let titleClass = '';
        
        if (msg.winner === 120) { // 'x' in ASCII
            title = 'Player X Won!';
            titleClass = 'x-wins';
        } else if (msg.winner === 111) { // 'o' in ASCII
            title = 'Player O Won!';
            titleClass = 'o-wins';
        } else {
            title = "It's a Tie!";
            titleClass = 'tie';
        }
        
        this.gameOverTitleEl.textContent = title;
        this.gameOverTitleEl.className = `game-over-title ${titleClass}`;
        this.gameOverOverlayEl.classList.remove('hidden');
        
        console.log('Game finished:', msg);
    }
    
    updateTimersFromServer(p1Time, p2Time, p1Move) {
        // Update timer values and sync with server
        this.p1Time = p1Time;
        this.p2Time = p2Time;
        this.currentPlayer = p1Move ? 1 : 2;
        this.lastUpdateTime = Date.now();
        
        // Update display immediately
        this.updateTimerDisplay();
        
        // Only start/restart the client-side timer if timers have been started
        if (this.timersStarted) {
            this.startTimer();
        }
    }
    
    startTimer() {
        // Clear existing timer
        this.stopTimer();
        
        // Start new timer that updates every 100ms for smooth countdown
        this.timerInterval = setInterval(() => {
            this.updateClientTimer();
        }, 100);
    }
    
    stopTimer() {
        if (this.timerInterval) {
            clearInterval(this.timerInterval);
            this.timerInterval = null;
        }
    }
    
    updateClientTimer() {
        const now = Date.now();
        const elapsed = now - this.lastUpdateTime;
        
        // Subtract elapsed time from current player's timer
        if (this.currentPlayer === 1) {
            this.p1Time = Math.max(0, this.p1Time - elapsed);
        } else {
            this.p2Time = Math.max(0, this.p2Time - elapsed);
        }
        
        this.lastUpdateTime = now;
        this.updateTimerDisplay();
        
        // Check if time ran out
        if ((this.currentPlayer === 1 && this.p1Time <= 0) || 
            (this.currentPlayer === 2 && this.p2Time <= 0)) {
            this.stopTimer();
        }
    }
    
    updateTimerDisplay() {
        // Only update timer display with actual values if timers have been started
        if (this.timersStarted) {
            this.p1TimerEl.textContent = this.formatTime(this.p1Time);
            this.p2TimerEl.textContent = this.formatTime(this.p2Time);
        }
        // Otherwise keep showing --:-- (set in handleOpponentFound)
    }
    
    formatTime(milliseconds) {
        const totalSeconds = Math.max(0, Math.ceil(milliseconds / 1000));
        const minutes = Math.floor(totalSeconds / 60);
        const seconds = totalSeconds % 60;
        return `${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
    }
    
    updateGameUI(p1Move) {
        // Update current turn indicator
        if (p1Move) {
            this.currentTurnEl.textContent = "Player X's Turn";
            this.player1InfoEl.className = "player-info active";
            this.player2InfoEl.className = "player-info inactive";
        } else {
            this.currentTurnEl.textContent = "Player O's Turn";
            this.player1InfoEl.className = "player-info inactive";
            this.player2InfoEl.className = "player-info active";
        }
    }
    
    highlightActiveBoard() {
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
    
    updateMiniboardWinners(completeBoards) {
        console.log('Raw complete_boards:', completeBoards);
        
        // Decode base64 string to get the actual byte values
        let decodedBytes;
        try {
            // Convert base64 string to byte array
            const binaryString = atob(completeBoards);
            decodedBytes = new Uint8Array(binaryString.length);
            for (let i = 0; i < binaryString.length; i++) {
                decodedBytes[i] = binaryString.charCodeAt(i);
            }
            console.log('Decoded bytes:', Array.from(decodedBytes));
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
                console.log(`Mini-board ${i}: byte value ${winnerByte}`);
                
                let winner = null;
                if (winnerByte === 120) { // 'x' is ASCII 120
                    winner = 'x';
                } else if (winnerByte === 111) { // 'o' is ASCII 111
                    winner = 'o';
                } else if (winnerByte === 0) {
                    // Board not completed, remove completed class
                    miniBoard.classList.remove('completed');
                    continue;
                }
                
                if (winner) {
                    console.log(`Mini-board ${i} won by: ${winner}`);
                    // Add completed class
                    miniBoard.classList.add('completed');
                    
                    // Create winner overlay
                    const winnerOverlay = document.createElement('div');
                    winnerOverlay.className = `mini-board-winner ${winner}-winner`;
                    winnerOverlay.textContent = winner.toUpperCase();
                    miniBoard.appendChild(winnerOverlay);
                } else {
                    // Remove completed class if winner is not set
                    miniBoard.classList.remove('completed');
                }
            } else {
                // Remove completed class if no data
                miniBoard.classList.remove('completed');
            }
        }
    }
    
    renderBoard(state) {
        if (state.length !== 81) return;
        
        for (let i = 0; i < 81; i++) {
            const cell = document.getElementById(`cell-${i}`);
            const value = state[i] || '';
            
            cell.textContent = value;
            
            // Add CSS classes for X and O
            cell.classList.remove('x', 'o');
            if (value.toLowerCase() === 'x') {
                cell.classList.add('x');
            } else if (value.toLowerCase() === 'o') {
                cell.classList.add('o');
            }
        }
    }
    
    handleCellClick(miniIndex, cellIndex) {
        console.log(`Cell clicked: mini-board ${miniIndex}, cell ${cellIndex}`);
        console.log('Current boardToPlayIn:', this.boardToPlayIn);
        
        // Don't allow clicks if timers haven't started yet
        if (!this.timersStarted) {
            console.log('Click blocked: game has not started yet');
            return;
        }
        
        // Don't allow clicks on completed mini-boards
        const miniBoard = document.getElementById(`mini-board-${miniIndex}`);
        if (miniBoard.classList.contains('completed')) {
            console.log('Click blocked: mini-board is completed');
            return;
        }
        
        // Check if click is in valid board (if specific board is highlighted)
        if (this.boardToPlayIn !== '?' && parseInt(this.boardToPlayIn) !== miniIndex) {
            console.log(`Click blocked: must play in board ${this.boardToPlayIn}, clicked on board ${miniIndex}`);
            return;
        }
        
        // Check if cell is already occupied
        const cell = document.getElementById(`cell-${miniIndex * 9 + cellIndex}`);
        if (cell.textContent.trim() !== '') {
            console.log('Click blocked: cell is already occupied');
            return;
        }
        
        // Send move to server
        const moveData = `${miniIndex}${cellIndex}`;
        console.log('Sending move to server:', moveData);
        this.socket.send(moveData);
    }
}

// Initialize the game when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    new GameBoard();
});
