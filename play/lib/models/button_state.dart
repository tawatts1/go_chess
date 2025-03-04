//contains whether the button is visible in the ui, and if it is enabled. 
class ButtonState {
  bool isVisible = false;
  bool isEnabled = false;
  ButtonState(this.isVisible, this.isEnabled);
  @override
  bool operator ==(Object other) =>
    other is ButtonState &&
    other.runtimeType == runtimeType &&
    other.isVisible == isVisible && other.isEnabled == isEnabled;
  @override
  int get hashCode => Object.hash(isVisible, isEnabled);
}