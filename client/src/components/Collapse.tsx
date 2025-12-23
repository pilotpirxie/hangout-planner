import { useState, type ReactNode } from "react";

export const Collapse = ({
  title,
  children,
}: {
  title: string;
  children: ReactNode;
}) => {
  const [isOpen, setIsOpen] = useState(false);

  const handleToggle = () => {
    setIsOpen(prev => !prev);
  };

  return (
    <div>
      <div
        className="cursor-pointer text-decoration-underline"
        onClick={handleToggle}><i className={isOpen ? "ri-arrow-down-line" : "ri-arrow-right-line"} />{title}</div>
      {isOpen ? <div>{children}</div> : null}
    </div>
  );
};