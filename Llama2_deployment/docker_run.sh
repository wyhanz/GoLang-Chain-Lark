docker run -itd \
    --gpus='"device=2"' \
    --shm-size=5g \
    -v /mnt/matrix/matrix/Llama-2-13b-chat-hf:/data \
    -p 38880:80 \
    ghcr.io/huggingface/text-generation-inference:0.9.4 \
    --model-id /data

