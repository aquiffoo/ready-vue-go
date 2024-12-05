<template>
	<div id="app">
		<h1>Ready, Vue, Go!</h1>

		<div v-if="!username">
			<input
				v-model="newUsername"
				type="text"
				placeholder="Enter username"
			/>
			<button @click="createUsername">Set Username</button>
		</div>

		<div v-if="username">
			<div class="chat-w" ref="chatWindow">
				<div
					v-for="(msg, index) in messages"
					:key="index"
					class="message"
				>
					<p><strong>{{ msg.sender }}:</strong> {{ msg.message }}</p>
				</div>
			</div>
			<input
				v-model="message"
				type="text"
				placeholder="Type a message..."
				@keyup.enter="sendMessage"
			/>
			<button @click="sendMessage">Send</button>
			<div v-if="typingUsers && typingUsers.length > 0" class="typing-indicator">
				<p v-for="user in typingUsers" :key="user">{{ user }} is typing...</p>
			</div>
		</div>
	</div>
</template>

<script>
export default {
	data() {
		return {
			username: null,
			newUsername: "",
			message: "",
			messages: [],
			isTyping: false,
		};
	},
	methods: {
		async createUsername() {
			if (this.newUsername.trim()) {
				try {
					const response = await fetch("http://localhost:8080/users", {
						method: "POST",
						headers: {
							"Content-Type": "application/json",
						},
						body: JSON.stringify({ username: this.newUsername }),
					});

					if (response.ok) {
						this.username = this.newUsername;
						this.fetchMessages();
					} else {
						console.error("Failed to create username");
					}
				} catch (error) {
					console.error("Error creating username:", error);
				}
			}
		},
		async sendMessage() {
			if (this.message.trim()) {
				try {
					const response = await fetch("http://localhost:8080/message", {
						method: "POST",
						headers: {
							"Content-Type": "application/json",
						},
						body: JSON.stringify({
							message: this.message,
							sender: this.username,
						}),
					});

					if (response.ok) {
						this.message = "";
						await this.fetchMessages();
						this.scrollToLastMessage();
						await this.updateTypingStatus(false);
					} else {
						console.error("Failed to send message");
					}
				} catch (error) {
					console.error("Error sending message:", error);
				}
			}
		},
		async fetchMessages() {
			try {
				const response = await fetch("http://localhost:8080/messages");
				const messages = await response.json();
				this.messages = messages;
				this.scrollToLastMessage();
			} catch (error) {
				console.error("Error fetching messages:", error);
			}
		},
		scrollToLastMessage() {
			this.$nextTick(() => {
				const chatWindow = this.$refs.chatWindow;
				if (chatWindow) {
					chatWindow.scrollTop = chatWindow.scrollHeight;
				}
			});
		},
	},
	watch: {
		message(newVal, oldVal) {
			if (newVal && !oldVal) {
				this.isTyping = true; 
			} else if (!newVal && oldVal) {
				this.isTyping = false;
			}
		},
	},
	mounted() {
		this.fetchMessages();
		setInterval(this.fetchMessages, 100);
		setInterval(this.fetchTypingUsers, 100);
	},
};
</script>

<style scoped>
@import url('https://fonts.googleapis.com/css2?family=Funnel+Sans:ital,wght@0,300..800;1,300..800&family=IBM+Plex+Mono:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;1,100;1,200;1,300;1,400;1,500;1,600;1,700&family=IBM+Plex+Sans:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;1,100;1,200;1,300;1,400;1,500;1,600;1,700&display=swap');

* {
	font-family: 'IBM Plex Sans', monospace;
	font-size: 2rem;
}

.chat-w {
	width: 500px;
	height: 700px;
	border: 1px solid #ccc;
	overflow-y: auto;
}

.message {
	margin-bottom: 10px;
}

.message p {
	margin: 0;
	padding: 10px;
	border-radius: 5px;
	background-color: #f2f2f2;
}

.message p strong {
	color: #333;
}

button {
	margin-left: 10px;
	padding: 5px 10px;
	border: none;
	border-radius: 5px;
	background-color: #007bff;
	color: #fff;
	cursor: pointer;
}

.typing-indicator p {
	margin: 0;
	padding: 5px;
	color: #555;
	font-style: italic;
	font-size: 1.5rem;
}


input {
    margin-bottom: 10px;
	padding: 5px;
}
</style>
