import styles from '@/styles/chatButton.module.scss';

interface ChatButtonProps {
  label: string;
  onClick: () => void;
  disabled?: boolean;
}

export default function ChatButton({ label, onClick, disabled }: ChatButtonProps) {
  return (
    <button className={styles.button} onClick={onClick} disabled={disabled}>
      {label}
    </button>
  );
}
