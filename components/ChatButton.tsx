import styles from '@/styles/chatButton.module.scss';

interface ChatButtonProps {
  label: string;
  onClick: () => void;
}

export default function ChatButton({ label, onClick }: ChatButtonProps) {
  return (
    <button className={styles.button} onClick={onClick}>
      {label}
    </button>
  );
}
