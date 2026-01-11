(function () {
  'use strict';

  // DOM Elements
  const form = document.getElementById('uploadForm');
  const dropzone = document.getElementById('dropzone');
  const fileInput = dropzone.querySelector('input[type="file"]');
  const fileInfo = document.getElementById('fileInfo');
  
  const progressSection = document.getElementById('progressSection');
  const progressBar = document.getElementById('progressBar');
  const progressPercent = document.querySelector('.progress-percent');
  
  const statusSection = document.getElementById('statusSection');
  const statusBadge = document.getElementById('statusBadge');
  const statusMessage = document.getElementById('statusMessage');
  
  const downloadSection = document.getElementById('downloadSection');
  const downloadBtn = document.getElementById('downloadBtn');
  
  const errorSection = document.getElementById('errorSection');
  const errorMessage = document.getElementById('errorMessage');

  // Utility functions
  function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  }

  function hideAllSections() {
    progressSection.style.display = 'none';
    statusSection.style.display = 'none';
    downloadSection.style.display = 'none';
    errorSection.style.display = 'none';
  }

  function showError(message) {
    hideAllSections();
    errorSection.style.display = 'block';
    errorMessage.textContent = message;
    form.classList.remove('loading');
  }

  function updateProgress(percent) {
    progressBar.style.width = percent + '%';
    progressPercent.textContent = percent + '%';
  }

  function updateStatus(status) {
    const statusConfig = {
      'PENDING': {
        class: 'pending',
        text: 'Pending',
        message: 'Your file is queued for processing...'
      },
      'RUNNING': {
        class: 'running',
        text: 'Processing',
        message: 'Optimizing your file...'
      },
      'DONE': {
        class: 'done',
        text: 'Complete',
        message: 'Your file has been optimized successfully!'
      },
      'FAILED': {
        class: 'failed',
        text: 'Failed',
        message: 'Something went wrong during processing'
      }
    };

    const config = statusConfig[status] || statusConfig['PENDING'];
    
    statusBadge.className = 'status-badge ' + config.class;
    statusBadge.textContent = config.text;
    statusMessage.textContent = config.message;
  }

  // File selection handlers
  function handleFileSelect(file) {
    if (!file) return;

    const fileName = file.name;
    const fileSize = formatFileSize(file.size);

    // Update file info display
    fileInfo.style.display = 'flex';
    fileInfo.querySelector('.file-name').textContent = fileName;
    fileInfo.querySelector('.file-size').textContent = fileSize;

    // Update icon based on file type
    const fileIcon = fileInfo.querySelector('.file-icon');
    if (file.type.startsWith('image/')) {
      fileIcon.textContent = '🖼️';
    } else if (file.type.startsWith('video/')) {
      fileIcon.textContent = '🎥';
    } else if (file.type === 'application/pdf') {
      fileIcon.textContent = '📄';
    } else {
      fileIcon.textContent = '📎';
    }
  }

  // Dropzone event handlers
  dropzone.addEventListener('click', (e) => {
    if (e.target !== fileInput) {
      fileInput.click();
    }
  });

  dropzone.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropzone.classList.add('dragover');
  });

  dropzone.addEventListener('dragleave', () => {
    dropzone.classList.remove('dragover');
  });

  dropzone.addEventListener('drop', (e) => {
    e.preventDefault();
    dropzone.classList.remove('dragover');
    
    if (e.dataTransfer.files.length > 0) {
      fileInput.files = e.dataTransfer.files;
      handleFileSelect(e.dataTransfer.files[0]);
    }
  });

  fileInput.addEventListener('change', (e) => {
    if (e.target.files.length > 0) {
      handleFileSelect(e.target.files[0]);
    }
  });

  // Form submission
  form.addEventListener('submit', async (e) => {
    e.preventDefault();
    
    // Reset UI
    hideAllSections();

    if (!fileInput.files.length) {
      showError('Please select a file first.');
      return;
    }

    // Show loading state
    form.classList.add('loading');

    try {
      // Upload file
      const formData = new FormData();
      formData.append('file', fileInput.files[0]);

      const uploadResponse = await fetch('/upload', {
        method: 'POST',
        body: formData,
      });

      if (!uploadResponse.ok) {
        const errorText = await uploadResponse.text();
        throw new Error(errorText || 'Upload failed. Please try again.');
      }

      const { job_id } = await uploadResponse.json();

      // Show progress
      progressSection.style.display = 'block';
      statusSection.style.display = 'block';
      updateProgress(0);
      updateStatus('PENDING');

      // Start polling
      pollJobStatus(job_id);

    } catch (error) {
      console.error('Upload error:', error);
      showError(error.message || 'Failed to upload file. Please try again.');
    }
  });

  // Job status polling
  async function pollJobStatus(jobId) {
    try {
      const response = await fetch(`/status/${jobId}`);

      if (!response.ok) {
        // Retry after delay
        setTimeout(() => pollJobStatus(jobId), 1000);
        return;
      }

      const job = await response.json();

      // Update progress
      updateProgress(job.Progress || 0);
      updateStatus(job.Status);

      // Handle completion
      if (job.Status === 'DONE') {
        form.classList.remove('loading');
        downloadSection.style.display = 'block';
        downloadBtn.href = `/download/${jobId}`;
        
        // Celebration animation
        confetti();
        return;
      }

      // Handle failure
      if (job.Status === 'FAILED') {
        showError(job.Error || 'Compression failed. Please try again.');
        return;
      }

      // Continue polling
      setTimeout(() => pollJobStatus(jobId), 1000);

    } catch (error) {
      console.error('Status polling error:', error);
      // Retry after delay
      setTimeout(() => pollJobStatus(jobId), 1000);
    }
  }

  // Simple confetti effect
  function confetti() {
    const duration = 2000;
    const end = Date.now() + duration;

    (function frame() {
      // Create confetti elements
      for (let i = 0; i < 3; i++) {
        const confetti = document.createElement('div');
        confetti.style.position = 'fixed';
        confetti.style.width = '10px';
        confetti.style.height = '10px';
        confetti.style.background = ['#667eea', '#764ba2', '#10b981', '#f59e0b'][Math.floor(Math.random() * 4)];
        confetti.style.left = Math.random() * window.innerWidth + 'px';
        confetti.style.top = '-10px';
        confetti.style.borderRadius = '50%';
        confetti.style.pointerEvents = 'none';
        confetti.style.zIndex = '9999';
        confetti.style.opacity = '1';
        confetti.style.transition = 'all 2s ease-out';
        
        document.body.appendChild(confetti);

        // Animate
        setTimeout(() => {
          confetti.style.top = window.innerHeight + 'px';
          confetti.style.opacity = '0';
        }, 10);

        // Remove after animation
        setTimeout(() => {
          confetti.remove();
        }, 2100);
      }

      if (Date.now() < end) {
        requestAnimationFrame(frame);
      }
    })();
  }

  // Reset form on page load
  window.addEventListener('load', () => {
    form.reset();
    fileInfo.style.display = 'none';
    hideAllSections();
  });

})();