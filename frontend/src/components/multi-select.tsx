import { useMemo } from 'react';
import { type MultiValue } from 'react-select';
import CreatableSelect from 'react-select/creatable';

interface Props {
  value?: string[];
  options?: { label: string; value: string }[];
  onChange: (value?: string[]) => void;
  onCreate: (value: string) => void;
  disabled?: boolean;
  placeholder?: string;
}

export function MultiSelect({
  options = [],
  value = [],
  onChange,
  onCreate,
  placeholder,
  disabled
}: Props) {
  const handleChange = (selected: MultiValue<{ label: string; value: string }>) => {
    const newValues = selected.map((option) => option.value);
    onChange?.(newValues);
  };

  const formattedValue = useMemo(() => {
    return options.filter((option) => value.includes(option.value));
  }, [options, value]);

  return (
    <CreatableSelect
      isMulti
      isDisabled={disabled}
      placeholder={placeholder}
      onChange={handleChange}
      onCreateOption={onCreate}
      value={formattedValue}
      options={options}
      className="text-sm  h-10"
      styles={{
        control: (base) => ({
          ...base,
          color: '#E5E7EB',
          borderRadius: '0.5rem',
          background: '#0F111A',
          borderColor: '#E2E8F0',
          ':hover': {
            borderColor: '#E2E8F0'
          }
        }),
        option: (provided) => ({
          ...provided,
          background: '#1A1C25',
          borderColor: '#E2E8F0',
          ':hover': {
            background: '#E5E7EB',
            color: '#1A1C25'
          }
        }),
        input: (styles) => ({
          ...styles,
          color: '#E5E7EB'
        }),
        placeholder: (styles) => ({
          ...styles,
          color: '#616975'
        }),
        multiValueLabel: (styles) => ({
          ...styles,
          background: '#1A1C25',
          color: '#E5E7EB'
        }),
        multiValue: (styles) => ({
          ...styles,
          background: '#1A1C25',
          color: '#E5E7EB',
          border: '1px solid #E5E7EB',
          borderRadius: '0.2rem'
        }),
        multiValueRemove: (styles) => ({
          ...styles,
          background: '#1A1C25',
          color: '#E5E7EB'
        })
      }}
    />
  );
}
