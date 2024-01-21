// Function to handle opening tabs
function openTab(tabId) {
    var tabs = document.getElementsByClassName('tab');
    for (var i = 0; i < tabs.length; i++) {
        tabs[i].style.display = 'none';
    }
    document.getElementById(tabId).style.display = 'block';
}

// Function to fetch and display total transactions
function fetchTotalTransactions() {
    fetch('/api/members') // Endpoint to get all members and their transaction totals
        .then(response => response.json())
        .then(data => {
            const table = document.getElementById('transactionsTable');
            // Clear previous data
            while (table.rows.length > 1) {
                table.deleteRow(1);
            }
            // Add new data
            data.forEach(item => {
                let row = table.insertRow();
                let cell1 = row.insertCell(0);
                let cell2 = row.insertCell(1);
                cell1.innerHTML = item.member_id;
                cell2.innerHTML = item.total_fee;
            });
        });
}

// Function to handle add member form submission
document.getElementById('addForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const name = document.getElementById('addName').value;
    fetch('/api/member', { // Endpoint to add a member
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: name })
    }).then(response => {
        if (response.ok) {
            alert('客戶成功新增');
            // Optionally refresh the page or update the UI
        } else {
            alert('新增客戶時出錯');
        }
    });
});

// Function to handle update member form submission
document.getElementById('modifyForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const id = document.getElementById('modifyId').value;
    const name = document.getElementById('modifyName').value;
    fetch('/api/member/' + id, { // Endpoint to update a member
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username: name })
    }).then(response => {
        if (response.ok) {
            alert('客戶資料已更新');
            // Optionally refresh the page or update the UI
        } else {
            alert('更新客戶資料時出錯');
        }
    });
});

// Function to handle search form submission
document.getElementById('searchForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const memberId = document.getElementById('searchId').value;
    const startDate = document.getElementById('startDate').value;
    const endDate = document.getElementById('endDate').value;

    // Clear previous search results
    const table = document.getElementById('searchResultsTable');
    while (table.rows.length > 1) {
        table.deleteRow(1);
    }

    // Fetch and display new search results
    fetch(`/api/member/${memberId}/transactions?start=${startDate}&end=${endDate}`) // Endpoint to get transactions of a member
        .then(response => response.json())
        .then(data => {
            data.forEach(item => {
                let row = table.insertRow();
                let cell1 = row.insertCell(0);
                let cell2 = row.insertCell(1);
                let cell3 = row.insertCell(2);
                cell1.innerHTML = item.id;
                cell2.innerHTML = item.borrow_fee;
                cell3.innerHTML = item.create_time;
            });
        });
});

// Initial fetch of total transactions
document.addEventListener('DOMContentLoaded', function() {
    fetchTotalTransactions();
});
